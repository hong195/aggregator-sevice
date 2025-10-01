package main

import (
	"log"

	"github.com/hong195/aggregator-sevice/config"
	"github.com/hong195/aggregator-sevice/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
