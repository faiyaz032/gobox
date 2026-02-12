package main

import (
	"log"

	"github.com/faiyaz032/gobox/cmd"
	"github.com/faiyaz032/gobox/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cmd.RunMigrate(cfg); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Start the server
	cmd.RunServer(cfg)
}
