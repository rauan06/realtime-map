package config

import (
	"errors"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/rauan06/realtime-map/analytics/internal/domain"
)

type (
	// Config -.
	Config struct {
		App App
		Log Log
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
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, domain.ErrConfigFileLode
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Join(domain.ErrConfigError, err)
	}

	return cfg, nil
}
