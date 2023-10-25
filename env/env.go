package env

import (
	"context"
	"fmt"

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

	Port        uint16 `env:"PORT,default=8000"`
	ServiceName string `env:"SERVICE_NAME,default=unknown"`
}

func LoadConfig(ctx context.Context, envPath string) (*Config, error) {
	var c Config

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %v", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
