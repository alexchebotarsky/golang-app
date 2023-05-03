package main

import (
	"context"
	"log"

	"github.com/goodleby/pure-go-server/config"
	"github.com/goodleby/pure-go-server/server"
)

func main() {
	ctx := context.Background()

	config, err := config.Load(ctx)
	if err != nil {
		log.Printf("error loading config: %v", err)
	}

	server, err := server.New(ctx, config)
	if err != nil {
		log.Printf("error creating server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Printf("error starting server: %v", err)
	}
}
