package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig represents the configuration for CORS middleware
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       int
	Credentials  bool
}

type CORSMiddleware interface {
	HandleCORS() gin.HandlerFunc
}

type CORSMiddlewareImpl struct {
	CORSConfig *CORSConfig
}

// Make sure that CORSMiddlewareImpl implements CORSMiddleware
var _ CORSMiddleware = (*CORSMiddlewareImpl)(nil)

func NewCORSMiddleware(cfg *config.ApiConfig) (*CORSMiddlewareImpl, error) {
	return &CORSMiddlewareImpl{
		CORSConfig: DefaultCORSConfig(cfg),
	}, nil
}

// DefaultCORSConfig returns the default CORS configuration
func DefaultCORSConfig(cfg *config.ApiConfig) *CORSConfig {
	return &CORSConfig{
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: cfg.AllowMethods,
		AllowHeaders: cfg.AllowHeaders,
		MaxAge:       cfg.AccessControlMaxAge,
		Credentials:  cfg.IsCredentials,
	}
}

// HandleCORS creates a middleware function to handle CORS
func (cm *CORSMiddlewareImpl) HandleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		isDebug := gin.Mode() == gin.DebugMode

		// Check whether origin is authorized
		if (isDebug && strings.HasPrefix(origin, "http://localhost:")) || contains(cm.CORSConfig.AllowOrigins, origin) {
			allowOrigin = origin
		}

		if allowOrigin != "" {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", joinStrings(cm.CORSConfig.AllowMethods))
			c.Header("Access-Control-Allow-Headers", joinStrings(cm.CORSConfig.AllowHeaders))
			c.Header("Access-Control-Max-Age", strconv.Itoa(cm.CORSConfig.MaxAge))

			if cm.CORSConfig.Credentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// joinStrings joins string slice with comma
func joinStrings(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for _, s := range strings[1:] {
		result += ", " + s
	}
	return result
}

// Check whether a value is present in a slice
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
