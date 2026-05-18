package pipeline

import (
	"context"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/extractor"
	"github.com/rauan06/realtime-map/etl/internal/loader"
	"github.com/rauan06/realtime-map/etl/internal/transformer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

type Pipeline struct {
	name        string
	extractor   extractor.Extractor
	transformer transformer.Transformer
	loader      loader.Loader
	l           logger.Interface

	fetchInterval time.Duration
	flushInterval time.Duration
	batchSize     int
}

type Config struct {
	Name          string
	Extractor     extractor.Extractor
	Transformer   transformer.Transformer
	Loader        loader.Loader
	Logger        logger.Interface
	FetchInterval time.Duration
	FlushInterval time.Duration
	BatchSize     int
}

func New(cfg Config) *Pipeline {
	return &Pipeline{
		name:          cfg.Name,
		extractor:     cfg.Extractor,
		transformer:   cfg.Transformer,
		loader:        cfg.Loader,
		l:             cfg.Logger,
		fetchInterval: cfg.FetchInterval,
		flushInterval: cfg.FlushInterval,
		batchSize:     cfg.BatchSize,
	}
}

func (p *Pipeline) Run(ctx context.Context) error {
	fetchTicker := time.NewTicker(p.fetchInterval)
	defer fetchTicker.Stop()

	flushTicker := time.NewTicker(p.flushInterval)
	defer flushTicker.Stop()

	p.l.Info("[%s] pipeline started: fetch every %s, flush every %s or %d msgs",
		p.name, p.fetchInterval, p.flushInterval, p.batchSize)

	for {
		select {
		case <-ctx.Done():
			p.l.Info("[%s] shutting down, flushing remaining events...", p.name)

			if err := p.loader.Flush(context.Background()); err != nil {
				p.l.Error("[%s] final flush error: %s", p.name, err)
			}

			return nil

		case <-fetchTicker.C:
			p.fetchAndBuffer(ctx)

		case <-flushTicker.C:
			if p.loader.Len() > 0 {
				p.l.Debug("[%s] flush ticker: flushing %d events", p.name, p.loader.Len())

				if err := p.loader.Flush(ctx); err != nil {
					p.l.Error("[%s] flush error: %s", p.name, err)
				}
			}
		}
	}
}

func (p *Pipeline) fetchAndBuffer(ctx context.Context) {
	records, err := p.extractor.Extract(ctx)
	if err != nil {
		p.l.Error("[%s] extract error: %s", p.name, err)

		return
	}

	if len(records) == 0 {
		p.l.Debug("[%s] no records extracted", p.name)

		return
	}

	events, err := p.transformer.Transform(records)
	if err != nil {
		p.l.Error("[%s] transform error: %s", p.name, err)

		return
	}

	for _, event := range events {
		p.loader.Add(event)
	}

	p.l.Info("[%s] buffered %d events (total: %d)", p.name, len(events), p.loader.Len())

	if p.loader.Len() >= p.batchSize {
		p.l.Info("[%s] batch size reached (%d >= %d), flushing", p.name, p.loader.Len(), p.batchSize)

		if err := p.loader.Flush(ctx); err != nil {
			p.l.Error("[%s] batch flush error: %s", p.name, err)
		}
	}
}
