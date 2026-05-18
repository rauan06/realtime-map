package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/rauan06/realtime-map/api-gateway/internal/domain"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		GRPC    GRPC
		Metrics Metrics
		Swagger Swagger
		Auth    Auth
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// GRPC -.
	GRPC struct {
		Port              string `env:"GRPC_PORT,required"`
		ReflectionEnabled bool   `env:"GRPC_REFLECTION_ENABLED" envDefault:"true"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// Auth gates the WebSocket / device endpoints behind HS256 JWTs. Disable
	// for local development by setting AUTH_ENABLED=false.
	Auth struct {
		Enabled      bool          `env:"AUTH_ENABLED"        envDefault:"false"`
		JWTSecret    string        `env:"AUTH_JWT_SECRET"     envDefault:""`
		TokenTTL     time.Duration `env:"AUTH_TOKEN_TTL"      envDefault:"24h"`
		SharedSecret string        `env:"AUTH_SHARED_SECRET"  envDefault:""`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.Join(err, domain.ErrConfigFileLode)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
