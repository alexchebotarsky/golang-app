package config

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port Port `env:"PORT,default=8000"`
}

type Port = uint16

func Load(ctx context.Context) (*Config, error) {
	var c Config

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("error loading environment variables: %v", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		return &c, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
