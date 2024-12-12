package main

import (
	"github.com/gin-gonic/gin"
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

	components := initialisation.InitializeAppComponents(db, cfg)
	middlewares, mcErr := initialisation.InitializeMiddlewareComponents(cfg)
	if mcErr != nil {
		log.Fatalf("Failed to initialize app components: %v", err)
	}
	// Gin application configuration
	r := gin.Default()
	initialisation.ConfigureRoutes(r, components, middlewares)

	runError := r.Run()
	if runError != nil {
		return
	}
}
