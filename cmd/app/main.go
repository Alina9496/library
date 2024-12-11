package main

import (
	"log"

	"github.com/Alina9496/library/config"
	"github.com/Alina9496/library/internal/app"
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
