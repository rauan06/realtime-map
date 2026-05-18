package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var errLoadEnv = errors.New("loading .env file")

type (
	Config struct {
		App   App
		HTTP  HTTP
		Log   Log
		Kafka Kafka
	}

	App struct {
		Name    string `env:"APP_NAME"    envDefault:"map-dashboard"`
		Version string `env:"APP_VERSION" envDefault:"0.1.0"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" envDefault:"8090"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"INFO"`
	}

	Kafka struct {
		BootstrapServers string   `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"broker:29092"`
		GroupID          string   `env:"KAFKA_GROUP_ID"          envDefault:"dashboard"`
		Topics           []string `env:"KAFKA_TOPICS"            envSeparator:","                envDefault:"etl_flights,etl_ships,etl_transport,etl_roads,obu_data,anomalies"`
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
