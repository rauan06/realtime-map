package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var errLoadEnv = errors.New("loading .env file")

type (
	Config struct {
		App      App
		Log      Log
		Kafka    Kafka
		Detector Detector
	}

	App struct {
		Name    string `env:"APP_NAME"    envDefault:"map-anomaly"`
		Version string `env:"APP_VERSION" envDefault:"0.1.0"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"INFO"`
	}

	Kafka struct {
		BootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"broker:29092"`
		GroupID          string `env:"KAFKA_GROUP_ID"          envDefault:"anomaly"`
		FlightTopic      string `env:"KAFKA_FLIGHT_TOPIC"      envDefault:"etl_flights"`
		ShipTopic        string `env:"KAFKA_SHIP_TOPIC"        envDefault:"etl_ships"`
		OutputTopic      string `env:"KAFKA_OUTPUT_TOPIC"      envDefault:"anomalies"`
	}

	Detector struct {
		// Warmup observations per layer before scoring begins.
		Warmup int `env:"DETECTOR_WARMUP" envDefault:"500"`

		// Score threshold in [0,1]. Higher = fewer, more confident alerts.
		Threshold float64 `env:"DETECTOR_THRESHOLD" envDefault:"0.75"`

		// Cooldown silences re-fires for the same (layer, source_id) for
		// this long after an alert.
		Cooldown time.Duration `env:"DETECTOR_COOLDOWN" envDefault:"5m"`

		// Rolling sample buffer capacity per layer.
		BufferCap int `env:"DETECTOR_BUFFER_CAP" envDefault:"2000"`

		// Isolation Forest tuning.
		NumTrees   int `env:"IFOREST_TREES"       envDefault:"100"`
		SampleSize int `env:"IFOREST_SAMPLE_SIZE" envDefault:"256"`
	}
)

func NewConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("%w: %w", errLoadEnv, err)
		}
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
