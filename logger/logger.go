// Package logger provides structured logging capabilities for the application
// using zap as the underlying logging library. It supports both development
// and production environments with appropriate configuration.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

// Environment defines the application environment for configuring the logger
type Environment string

const (
	// DEVELOPMENT enables verbose logging with debug level information
	DEVELOPMENT Environment = "development"
	// PRODUCTION enables optimized logging for production use cases
	PRODUCTION Environment = "production"
)

var (
	// logger is the singleton zap logger instance
	logger *zap.Logger
	// sugarLogger is the singleton sugared logger instance for printf-style logging
	sugarLogger *zap.SugaredLogger
	// once ensures the logger is initialized only once
	once sync.Once
)

// Init initializes the global zap logger instance
// It creates a production or development logger based on the provided environment.
// This function is safe to call multiple times as it ensures singleton initialization.
//
// Parameters:
//   - environment: Environment - PRODUCTION or DEVELOPMENT to determine logging configuration
//
// Returns:
//   - *zap.Logger: The initialized logger instance
func Init(environment Environment) *zap.Logger {
	once.Do(func() {
		var err error

		// Define logger configuration based on environment
		var config zap.Config

		if environment == PRODUCTION {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
		}

		// Customize encoder config for better readability
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// Create the logger
		logger, err = config.Build()
		if err != nil {
			log.Fatalf("Failed to initialize zap logger: %v", err)
		}

		// Create sugared logger for convenience
		sugarLogger = logger.Sugar()

		// Replace the standard library's global logger
		redirectStdLogToZap()
	})

	return logger
}

// Get returns the global zap logger instance
// If the logger hasn't been initialized yet, it initializes it with development settings.
//
// Returns:
//   - *zap.Logger: The global logger instance
func Get() *zap.Logger {
	if logger == nil {
		// Default to development logger if Init wasn't called
		Init(DEVELOPMENT)
	}
	return logger
}

// redirectStdLogToZap redirects standard library's log to zap
// This ensures that all logging, even from dependencies using the standard library,
// goes through our configured zap logger for consistent output.
func redirectStdLogToZap() {
	// Redirect standard library's log to zap
	_ = func(msg string) {
		sugarLogger.Info(msg)
	}

	// Replace the standard library's global logger
	log.SetFlags(0) // Remove timestamp from standard logger output
	log.SetOutput(zap.NewStdLog(logger).Writer())
}
