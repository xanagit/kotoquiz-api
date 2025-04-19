package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/initialisation"
	"github.com/xanagit/kotoquiz-api/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.Init(logger.PRODUCTION)
	defer func() {
		// Handle potential error from sync operation
		if err := log.Sync(); err != nil {
			log.Fatal("Could not initialize the logger")
		}
	}()

	// Load configuration
	config.SetLogger(log)
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Failed to load config",
			zap.Error(err))
	}

	db, dbErr := initialisation.DatabaseConnectionFromConfig(cfg, log)
	if dbErr != nil {
		log.Fatal("Unable to connect to database",
			zap.Error(dbErr))
	}

	components := initialisation.InitializeAppComponents(db, cfg)
	middlewares, mcErr := initialisation.InitializeMiddlewareComponents(cfg, log)
	if mcErr != nil {
		log.Fatal("Failed to initialize app components",
			zap.Error(mcErr))
	}

	// Gin application configuration
	r := gin.Default()
	initialisation.ConfigureRoutes(r, components, middlewares, log)

	// Start the server
	log.Info("Starting server on port " + cfg.App.Port)
	runError := r.Run()
	if runError != nil {
		log.Error("Server failed to start",
			zap.Error(runError))
	}
}
