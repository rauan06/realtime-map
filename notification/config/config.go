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
		App      App
		Log      Log
		Kafka    Kafka
		Geofence Geofence
		Notifier Notifier
	}

	App struct {
		Name    string `env:"APP_NAME"    envDefault:"map-notification"`
		Version string `env:"APP_VERSION" envDefault:"0.1.0"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"INFO"`
	}

	Kafka struct {
		BootstrapServers string   `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"broker:29092"`
		GroupID          string   `env:"KAFKA_GROUP_ID"          envDefault:"notification"`
		Topics           []string `env:"KAFKA_TOPICS"            envSeparator:","                envDefault:"etl_flights,etl_ships,etl_transport,obu_data"`
	}

	Geofence struct {
		ConfigPath string `env:"GEOFENCE_CONFIG_PATH" envDefault:"geofences.json"`
	}

	Notifier struct {
		WebhookURL string `env:"WEBHOOK_URL" envDefault:""`
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
