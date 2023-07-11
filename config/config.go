package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config contains all environment configuration.
type Config struct {
	DatabaseUser     string `env:"DATABASE_USER,required"`
	DatabasePassword string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost     string `env:"DATABASE_HOST,default=localhost"`
	DatabasePort     Port   `env:"DATABASE_PORT,default=5432"`
	DatabaseName     string `env:"DATABASE_NAME,required"`
	DatabaseOptions  string `env:"DATABASE_OPTIONS,default=?sslmode=disable"`

	AuthSecret    string `env:"AUTH_SECRET,required"`
	AuthAdminKey  string `env:"AUTH_ADMIN_KEY,required"`
	AuthEditorKey string `env:"AUTH_EDITOR_KEY,required"`
	AuthViewerKey string `env:"AUTH_VIEWER_KEY,required"`

	Port Port `env:"PORT,default=8000"`
}

// Port is a valid port.
type Port = uint16

// Load loads the environment configuration.
func Load(ctx context.Context, envPath string) (*Config, error) {
	var c Config

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %v", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
