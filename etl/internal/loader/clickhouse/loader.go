package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	chgo "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const (
	dialTimeout     = 10 * time.Second
	maxOpenConns    = 10
	maxIdleConns    = 5
	connMaxLifetime = time.Hour
	initialBuffer   = 512
)

// Loader buffers events and bulk-inserts them into the ClickHouse table
// `etl_events` (see migrations/clickhouse/001_init.sql). The table is a
// generic event sink with JSON payload — the dashboard reads it for
// historical replay and aggregate queries.
type Loader struct {
	conn   driver.Conn
	source string
	mu     sync.Mutex
	buffer []domain.KafkaEvent
}

type Options struct {
	Addr     string
	Database string
	Username string
	Password string
	Source   string // tag stored alongside each event ("flight" / "ship" / ...)
}

func New(ctx context.Context, opts Options) (*Loader, error) {
	conn, err := chgo.Open(&chgo.Options{
		Addr: []string{opts.Addr},
		Auth: chgo.Auth{
			Database: opts.Database,
			Username: opts.Username,
			Password: opts.Password,
		},
		DialTimeout:     dialTimeout,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse open: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close()

		return nil, fmt.Errorf("clickhouse ping: %w", err)
	}

	return &Loader{
		conn:   conn,
		source: opts.Source,
		buffer: make([]domain.KafkaEvent, 0, initialBuffer),
	}, nil
}

func (l *Loader) Close() error { return l.conn.Close() }

func (l *Loader) Add(event domain.KafkaEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buffer = append(l.buffer, event)
}

func (l *Loader) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.buffer)
}

func (l *Loader) Flush(ctx context.Context) error {
	l.mu.Lock()
	batch := make([]domain.KafkaEvent, len(l.buffer))
	copy(batch, l.buffer)
	l.buffer = l.buffer[:0]
	l.mu.Unlock()

	if len(batch) == 0 {
		return nil
	}

	bw, err := l.conn.PrepareBatch(ctx, "INSERT INTO etl_events (source, key, payload, received_at)")
	if err != nil {
		return fmt.Errorf("clickhouse prepare batch: %w", err)
	}

	now := time.Now().UTC()

	for _, ev := range batch {
		payload, err := json.Marshal(ev.Data)
		if err != nil {
			return fmt.Errorf("clickhouse marshal event %s: %w", ev.Key, err)
		}

		if err := bw.Append(l.source, ev.Key, string(payload), now); err != nil {
			return fmt.Errorf("clickhouse append: %w", err)
		}
	}

	if err := bw.Send(); err != nil {
		return fmt.Errorf("clickhouse send batch: %w", err)
	}

	return nil
}
