package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"sync/atomic"
)

var (
	// Single instance of the configuration
	instance atomic.Value
	// Protection for initialization
	once sync.Once
)

// Config Define config struct to hold the app configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Jwt      JwtConfig
	Cors     CorsConfig `yaml:"cors"`
}

// AppConfig and DatabaseConfig structs to structure config fields
type AppConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

type JwtConfig struct {
	SecretKey string
}

type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	MaxAge         int      `yaml:"max_age"`
}

// GetConfig returns the singleton instance of the configuration
// Thread-safe thanks to sync.Once and atomic.Value
func GetConfig() (*Config, error) {
	var loadErr error
	config := instance.Load()
	if config == nil {
		// If config is not loaded, load it in a thread-safe way
		once.Do(func() {
			cfg, err := loadConfig()
			if err != nil {
				loadErr = err
				return
			}
			if len(cfg.Jwt.SecretKey) < 32 {
				loadErr = fmt.Errorf("JWT secret key must be at least 32 bytes long")
			}
			instance.Store(cfg)
		})
		if loadErr != nil {
			return nil, loadErr
		}
		config = instance.Load()
	}
	return config.(*Config), nil
}

// loadConfig function loads and returns the configuration from config.yml and environment variables
// Creates a dedicated Viper instance to avoid conflicts
func loadConfig() (*Config, error) {
	v := viper.New() // Dedicated Viper instance
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// Enable environment variables override and prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// Bind environment variables for specific configuration fields
	envVars := map[string]string{
		"database.host":     "APP_DATABASE_HOST",
		"database.user":     "APP_DATABASE_USER",
		"database.password": "APP_DATABASE_PASSWORD",
		"database.name":     "APP_DATABASE_NAME",
		"database.port":     "APP_DATABASE_PORT",
		"jwt.secret-key":    "APP_JWT_SECRET_KEY",
	}
	for key, env := range envVars {
		if err := v.BindEnv(key, env); err != nil {
			log.Fatalf("Erreur lors de la liaison de '%s': %v", key, err)
		}
	}

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the configuration into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
