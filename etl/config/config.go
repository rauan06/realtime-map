package config

import (
	"errors"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

type (
	Config struct {
		App        App
		Log        Log
		Kafka      Kafka
		ClickHouse ClickHouse
		Sources    Sources
	}

	App struct {
		Name    string `env:"APP_NAME"     envDefault:"map-etl"`
		Version string `env:"APP_VERSION"  envDefault:"0.1.0"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"INFO"`
	}

	Kafka struct {
		BootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"broker:29092"`
		Enabled          bool   `env:"KAFKA_ENABLED"           envDefault:"true"`
	}

	ClickHouse struct {
		Addr     string `env:"CLICKHOUSE_ADDR"     envDefault:"clickhouse:9000"`
		Database string `env:"CLICKHOUSE_DB"       envDefault:"realtimedb"`
		Username string `env:"CLICKHOUSE_USER"     envDefault:"default"`
		Password string `env:"CLICKHOUSE_PASSWORD" envDefault:"example"`
		Enabled  bool   `env:"CLICKHOUSE_ENABLED"  envDefault:"true"`
	}

	// Sources controls per-source extractor cadence and topic mapping. Set
	// the *_ENABLED flag to false to skip a source entirely.
	Sources struct {
		Flight    SourceConfig `envPrefix:"FLIGHT_"`
		Ship      SourceConfig `envPrefix:"SHIP_"`
		Transport SourceConfig `envPrefix:"TRANSPORT_"`
		Road      SourceConfig `envPrefix:"ROAD_"`

		// Common pipeline tunables.
		FlushInterval time.Duration `env:"PIPELINE_FLUSH_INTERVAL" envDefault:"15s"`
		BatchSize     int           `env:"PIPELINE_BATCH_SIZE"     envDefault:"500"`
		HTTPTimeout   time.Duration `env:"HTTP_TIMEOUT"            envDefault:"15s"`
	}

	SourceConfig struct {
		Enabled       bool          `env:"ENABLED"        envDefault:"true"`
		Topic         string        `env:"TOPIC"`
		FetchInterval time.Duration `env:"FETCH_INTERVAL" envDefault:"30s"`
		// Endpoint is optional override of the default API base URL.
		Endpoint string `env:"ENDPOINT"`
	}
)

func NewConfig() (*Config, error) {
	// .env is optional: load it only if present so docker-compose env still works.
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, errors.Join(domain.ErrConfigFileLoad, err)
		}
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Join(domain.ErrConfigParse, err)
	}

	// Apply sane defaults for topics if unset.
	if cfg.Sources.Flight.Topic == "" {
		cfg.Sources.Flight.Topic = "etl_flights"
	}

	if cfg.Sources.Ship.Topic == "" {
		cfg.Sources.Ship.Topic = "etl_ships"
	}

	if cfg.Sources.Transport.Topic == "" {
		cfg.Sources.Transport.Topic = "etl_transport"
	}

	if cfg.Sources.Road.Topic == "" {
		cfg.Sources.Road.Topic = "etl_roads"
	}

	return cfg, nil
}
