package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App   App
		Log   Log
		GRPC  GRPC
		Kafka Kafka
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// Kafka -.
	Kafka struct {
		BootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost"`
		Topic            string `env:"KAFKA_TOPIC,required"`
	}

	// GRPC -.
	GRPC struct {
		Port              string `env:"GRPC_PORT,required"`
		ReflectionEnabled bool   `env:"GRPC_REFLECTION_ENABLED" envDefault:"true"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("loading .env %s", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
