package config

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config contains all environment configuration.
type Config struct {
	Port Port `env:"PORT,required"`
}

// Port is a valid port.
type Port = uint16

// Load loads the environment configuration.
func Load(ctx context.Context, envPath string) (*Config, error) {
	var c Config

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("error loading environment variables: %v", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		return &c, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
