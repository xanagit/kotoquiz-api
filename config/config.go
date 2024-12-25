package config

import (
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

// Define Config struct to hold the app configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
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

type AuthConfig struct {
	Keycloak KeycloakConfig
}

type KeycloakConfig struct {
	BaseUrl          string `mapstructure:"baseUrl"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	Realm            string `mapstructure:"realm"`
	AdminCliClientId string `mapstructure:"adminCliClientId"`
	ClientID         string `mapstructure:"clientId"`
	ClientSecret     string `mapstructure:"clientSecret"`
	RedirectURL      string `mapstructure:"redirectUrl"`
	IssuerURL        string `mapstructure:"issuerUrl"`
	CallbackURL      string `mapstructure:"callbackUrl"`
	LogoutURL        string `mapstructure:"logoutUrl"`
	CookieDomain     string `mapstructure:"cookieDomain"`
	CookieSecure     bool   `mapstructure:"cookieSecure"`
	CookieMaxAge     int    `mapstructure:"cookieMaxAge"`
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
			instance.Store(cfg)
		})
		if loadErr != nil {
			return nil, loadErr
		}
		config = instance.Load()
	}
	return config.(*Config), nil
}

// LoadConfig function loads and returns the configuration from config.yml and environment variables
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
		"database.host":                  "APP_DATABASE_HOST",
		"database.user":                  "APP_DATABASE_USER",
		"database.password":              "APP_DATABASE_PASSWORD",
		"database.name":                  "APP_DATABASE_NAME",
		"database.port":                  "APP_DATABASE_PORT",
		"auth.keycloak.baseUrl":          "APP_KEYCLOAK_BASE_URL",
		"auth.keycloak.user":             "APP_KEYCLOAK_USER",
		"auth.keycloak.password":         "APP_KEYCLOAK_PASSWORD",
		"auth.keycloak.realm":            "APP_KEYCLOAK_REALM",
		"auth.keycloak.adminCliClientId": "APP_KEYCLOAK_ADMIN_CLI_CLIENT_ID",
		"auth.keycloak.clientId":         "APP_KEYCLOAK_CLIENT_ID",
		"auth.keycloak.clientSecret":     "APP_KEYCLOAK_CLIENT_SECRET",
		"auth.keycloak.issuerUrl":        "APP_KEYCLOAK_ISSUER_URL",
		"auth.keycloak.redirectUrl":      "APP_KEYCLOAK_REDIRECT_URL",
		"auth.keycloak.callbackUrl":      "APP_KEYCLOAK_CALLBACK_URL",
		"auth.keycloak.logoutUrl":        "APP_KEYCLOAK_LOGOUT_URL",
		"auth.keycloak.cookieDomain":     "APP_KEYCLOAK_COOKIE_DOMAIN",
		"auth.keycloak.cookieSecure":     "APP_KEYCLOAK_COOKIE_SECURE",
		"auth.keycloak.cookieMaxAge":     "APP_KEYCLOAK_COOKIE_MAX_AGE",
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
