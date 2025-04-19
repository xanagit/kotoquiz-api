package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

type Environment string

const (
	DEVELOPMENT Environment = "development"
	PRODUCTION  Environment = "production"
)

var (
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
	once        sync.Once
)

// Init initializes the global zap logger instance
// It creates a production or development logger based on the provided environment
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
func Get() *zap.Logger {
	if logger == nil {
		// Default to development logger if Init wasn't called
		Init(DEVELOPMENT)
	}
	return logger
}

// GetSugar returns the global sugared logger instance
//func GetSugar() *zap.SugaredLogger {
//	if sugarLogger == nil {
//		// Default to development logger if Init wasn't called
//		Init("development")
//	}
//	return sugarLogger
//}

// redirectStdLogToZap redirects standard library's log to zap
func redirectStdLogToZap() {
	// Redirect standard library's log to zap
	_ = func(msg string) {
		sugarLogger.Info(msg)
	}

	// Replace the standard library's global logger
	log.SetFlags(0) // Remove timestamp from standard logger output
	log.SetOutput(zap.NewStdLog(logger).Writer())
}
