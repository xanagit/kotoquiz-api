// Package main is the entry point for the KotoQuiz API application
// It initializes all required components and starts the HTTP server
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/initialisation"
	"github.com/xanagit/kotoquiz-api/logger"
	"go.uber.org/zap"
)

// main is the application entry point that performs the following:
// 1. Initializes the logger with the appropriate environment
// 2. Loads application configuration
// 3. Establishes a database connection
// 4. Initializes application components and middleware
// 5. Configures routes and handlers
// 6. Starts the HTTP server
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

	// Connect to database
	db, dbErr := initialisation.DatabaseConnectionFromConfig(cfg, log)
	if dbErr != nil {
		log.Fatal("Unable to connect to database",
			zap.Error(dbErr))
	}

	// Initialize application components (repositories, services, controllers)
	components := initialisation.InitializeAppComponents(db, cfg)

	// Initialize middleware components (auth, CORS)
	middlewares, mcErr := initialisation.InitializeMiddlewareComponents(cfg, log)
	if mcErr != nil {
		log.Fatal("Failed to initialize app components",
			zap.Error(mcErr))
	}

	// Configure Gin application and routes
	r := gin.Default()
	initialisation.ConfigureRoutes(r, components, middlewares, log)

	// Start the HTTP server
	log.Info("Starting server on port " + cfg.App.Port)
	runError := r.Run()
	if runError != nil {
		log.Error("Server failed to start",
			zap.Error(runError))
	}
}
