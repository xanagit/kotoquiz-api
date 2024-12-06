package main

import (
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/initialisation"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, dbErr := initialisation.DatabaseConnectionFromConfig(cfg)
	if dbErr != nil {
		log.Fatalf("Unabled to connect to database %v", err)
	}

	r := initialisation.GinHandlers(cfg, db)

	runError := r.Run()
	if runError != nil {
		return
	}
}
