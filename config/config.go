package config

import (
	"github.com/spf13/viper"
	"log"
)

// Define Config struct to hold the app configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
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

// LoadConfig function loads and returns the configuration from config.yml and environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Enable environment variables override and prefix
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// Bind environment variables for specific configuration fields
	envVars := map[string]string{
		"database.host":     "APP_DATABASE_HOST",
		"database.user":     "APP_DATABASE_USER",
		"database.password": "APP_DATABASE_PASSWORD",
		"database.name":     "APP_DATABASE_NAME",
		"database.port":     "APP_DATABASE_PORT",
	}
	for key, env := range envVars {
		if err := viper.BindEnv(key, env); err != nil {
			log.Fatalf("Erreur lors de la liaison de '%s': %v", key, err)
		}
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the configuration into the Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
