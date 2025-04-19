// Package config provides configuration management for the application
// It handles loading and parsing configuration from files and environment variables,
// and provides a thread-safe singleton access pattern.
package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

var (
	// Single instance of the configuration
	instance atomic.Value
	// Protection for initialization
	once sync.Once
	// Logger instance
	logger *zap.Logger
)

// SetLogger sets the logger for the config package
//
// Parameters:
//   - l: *zap.Logger - The logger instance to use for config operations
func SetLogger(l *zap.Logger) {
	logger = l
}

// Config is the main configuration struct that holds all application settings
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// AppConfig contains general application settings
type AppConfig struct {
	// Port is the HTTP server port to listen on
	Port string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	// Host is the database server hostname or IP address
	Host string
	// User is the database username
	User string
	// Password is the database user's password
	Password string
	// Name is the database name to connect to
	Name string
	// Port is the database server port
	Port int
}

// AuthConfig contains authentication and authorization settings
type AuthConfig struct {
	// Keycloak contains Keycloak authentication provider settings
	Keycloak KeycloakConfig `mapstructure:"keycloak"`
	// ApiConfig contains API security settings like CORS
	ApiConfig ApiConfig `mapstructure:"apiConfig"`
}

// KeycloakConfig contains settings for Keycloak authentication provider
type KeycloakConfig struct {
	// BaseUrl is the base URL for Keycloak server
	BaseUrl string `mapstructure:"baseUrl"`
	// User is the Keycloak admin username
	User string `mapstructure:"user"`
	// Password is the Keycloak admin password
	Password string `mapstructure:"password"`
	// Realm is the Keycloak realm name for the application
	Realm string `mapstructure:"realm"`
	// AdminCliClientId is the client ID for admin operations
	AdminCliClientId string `mapstructure:"adminCliClientId"`
	// ClientID is the OAuth client ID for the application
	ClientID string `mapstructure:"clientId"`
	// ClientSecret is the OAuth client secret
	ClientSecret string `mapstructure:"clientSecret"`
	// RedirectURL is the OAuth redirect URL after login
	RedirectURL string `mapstructure:"redirectUrl"`
	// IssuerURL is the OpenID Connect issuer URL
	IssuerURL string `mapstructure:"issuerUrl"`
	// CallbackURL is the OpenID Connect callback URL
	CallbackURL string `mapstructure:"callbackUrl"`
	// LogoutURL is the OpenID Connect logout URL
	LogoutURL string `mapstructure:"logoutUrl"`
	// CookieDomain is the domain for authentication cookies
	CookieDomain string `mapstructure:"cookieDomain"`
	// CookieSecure indicates whether cookies should be sent only over HTTPS
	CookieSecure bool `mapstructure:"cookieSecure"`
	// CookieMaxAge is the authentication cookie expiration time in seconds
	CookieMaxAge int `mapstructure:"cookieMaxAge"`
}

// ApiConfig contains API security configuration
type ApiConfig struct {
	// AllowOrigins is a list of allowed CORS origins
	AllowOrigins []string `mapstructure:"allowOrigins"`
	// AllowMethods is a list of allowed HTTP methods for CORS
	AllowMethods []string `mapstructure:"allowMethods"`
	// AllowHeaders is a list of allowed HTTP headers for CORS
	AllowHeaders []string `mapstructure:"allowHeaders"`
	// AccessControlMaxAge is the max age in seconds for CORS preflight responses
	AccessControlMaxAge int `mapstructure:"accessControlMaxAge"`
	// IsCredentials enables support for cookies in CORS requests
	IsCredentials bool `mapstructure:"isCredentials"`
}

// GetConfig returns the singleton instance of the configuration
// Thread-safe thanks to sync.Once and atomic.Value
//
// Returns:
//   - *Config: The application configuration
//   - error: An error if configuration loading fails
func GetConfig() (*Config, error) {
	if logger == nil {
		// Create a default logger if not set
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}

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

// loadConfig loads and returns the configuration from config.yml and environment variables
// Creates a dedicated Viper instance to avoid conflicts
//
// Returns:
//   - *Config: The loaded configuration
//   - error: An error if configuration loading fails
func loadConfig() (*Config, error) {
	if logger == nil {
		// Create a default logger if not set
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}

	v := viper.New() // Dedicated Viper instance
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// Enable environment variables override and prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// Bind environment variables for specific configuration fields
	envVars := map[string]string{
		"database.host":                      "APP_DATABASE_HOST",
		"database.user":                      "APP_DATABASE_USER",
		"database.password":                  "APP_DATABASE_PASSWORD",
		"database.name":                      "APP_DATABASE_NAME",
		"database.port":                      "APP_DATABASE_PORT",
		"auth.keycloak.baseUrl":              "APP_KEYCLOAK_BASE_URL",
		"auth.keycloak.user":                 "APP_KEYCLOAK_USER",
		"auth.keycloak.password":             "APP_KEYCLOAK_PASSWORD",
		"auth.keycloak.realm":                "APP_KEYCLOAK_REALM",
		"auth.keycloak.adminCliClientId":     "APP_KEYCLOAK_ADMIN_CLI_CLIENT_ID",
		"auth.keycloak.clientId":             "APP_KEYCLOAK_CLIENT_ID",
		"auth.keycloak.clientSecret":         "APP_KEYCLOAK_CLIENT_SECRET",
		"auth.keycloak.issuerUrl":            "APP_KEYCLOAK_ISSUER_URL",
		"auth.keycloak.redirectUrl":          "APP_KEYCLOAK_REDIRECT_URL",
		"auth.keycloak.callbackUrl":          "APP_KEYCLOAK_CALLBACK_URL",
		"auth.keycloak.logoutUrl":            "APP_KEYCLOAK_LOGOUT_URL",
		"auth.keycloak.cookieDomain":         "APP_KEYCLOAK_COOKIE_DOMAIN",
		"auth.keycloak.cookieSecure":         "APP_KEYCLOAK_COOKIE_SECURE",
		"auth.keycloak.cookieMaxAge":         "APP_KEYCLOAK_COOKIE_MAX_AGE",
		"auth.apiConfig.allowOrigins":        "APP_API_CONFIG_ALLOW_ORIGIN",
		"auth.apiConfig.allowMethods":        "APP_API_CONFIG_ALLOW_METHODS",
		"auth.apiConfig.allowHeaders":        "APP_API_CONFIG_ALLOW_HEADERS",
		"auth.apiConfig.accessControlMaxAge": "APP_API_CONFIG_ACCESS_CONTROL_MAX_AGE",
		"auth.apiConfig.isCredentials":       "APP_API_CONFIG_IS_CREDENTIAL",
	}
	for key, env := range envVars {
		if err := v.BindEnv(key, env); err != nil {
			logger.Error("Error binding environment variable",
				zap.String("key", key),
				zap.Error(err))
		}
	}

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		logger.Error("Error reading config file", zap.Error(err))
		return nil, err
	}

	// Unmarshal the configuration into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		logger.Error("Error unmarshalling config", zap.Error(err))
		return nil, err
	}

	logger.Info("Configuration loaded successfully")
	return &config, nil
}
