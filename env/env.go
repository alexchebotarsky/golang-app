package env

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUser     string `env:"DATABASE_USER,required"`
	DatabasePassword string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost     string `env:"DATABASE_HOST,default=localhost"`
	DatabasePort     uint16 `env:"DATABASE_PORT,default=5432"`
	DatabaseName     string `env:"DATABASE_NAME,required"`
	DatabaseOptions  string `env:"DATABASE_OPTIONS,default=?sslmode=disable"`

	AuthSecret    string `env:"AUTH_SECRET,required"`
	AuthAdminKey  string `env:"AUTH_ADMIN_KEY,required"`
	AuthEditorKey string `env:"AUTH_EDITOR_KEY,required"`
	AuthViewerKey string `env:"AUTH_VIEWER_KEY,required"`

	ExampleEndpoint string `env:"EXAMPLE_ENDPOINT,default=https://swapi.dev/api/people/1"`

	PubSubProjectID string `env:"PUBSUB_PROJECT_ID,required"`

	Port          uint16 `env:"PORT,default=8000"`
	AllowedOrigin string `env:"ALLOWED_ORIGIN,default="`
	ServiceName   string `env:"SERVICE_NAME,default=unknown"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var c Config

	// We are loading env variables from .env file only for local development
	if err := godotenv.Load(".env"); err != nil {
		slog.Debug("error loading .env file: %v", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
