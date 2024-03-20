package env

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ServiceName       string     `env:"SERVICE_NAME,default=unknown"`
	Environment       string     `env:"ENVIRONMENT,default=stage"`
	LogLevel          slog.Level `env:"LOG_LEVEL,default=debug"`
	TracingSampleRate float64    `env:"TRACING_SAMPLE_RATE,default=1"`

	Port          uint16 `env:"PORT,default=8000"`
	AllowedOrigin string `env:"ALLOWED_ORIGIN,default=http://localhost:3000"`

	DatabaseUser     string `env:"DATABASE_USER,required"`
	DatabasePassword string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost     string `env:"DATABASE_HOST,default=127.0.0.1"`
	DatabasePort     uint16 `env:"DATABASE_PORT,default=5432"`
	DatabaseName     string `env:"DATABASE_NAME,required"`
	DatabaseOptions  string `env:"DATABASE_OPTIONS,default=?sslmode=disable"`

	AuthSecret    string        `env:"AUTH_SECRET,required"`
	AuthTokenTTL  time.Duration `env:"AUTH_TOKEN_TTL,default=20m"`
	AuthAdminKey  string        `env:"AUTH_ADMIN_KEY,required"`
	AuthEditorKey string        `env:"AUTH_EDITOR_KEY,required"`
	AuthViewerKey string        `env:"AUTH_VIEWER_KEY,required"`

	GoogleApplicationCredentials string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
	PubSubProjectID              string `env:"PUBSUB_PROJECT_ID,required"`
	TraceProjectID               string `env:"TRACE_PROJECT_ID,required"`

	ExampleEndpoint string `env:"EXAMPLE_ENDPOINT,required"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var c Config

	// We are loading env variables from .env file only for local development
	err := godotenv.Load(".env")
	if err != nil {
		slog.Debug("error loading .env file: %v", err)
	}

	err = envconfig.Process(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
