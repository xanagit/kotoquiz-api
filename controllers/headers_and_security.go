package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HeadersAndSecurityController interface {
	AddHeaders(c *gin.Context)
	CspReport(c *gin.Context)
}

type HeadersAndSecurityControllerImpl struct {
	Config config.Config
}

func (hs *HeadersAndSecurityControllerImpl) AddHeaders(c *gin.Context) {
	allowedOrigins := hs.Config.Cors.AllowedOrigins

	// Get the request origin
	origin := c.Request.Header.Get("Origin")
	// Check if the origin is allowed
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			break
		}
	}
	// Headers needed for CORS preflight requests
	c.Header("Access-Control-Allow-Methods", strings.Join(hs.Config.Cors.AllowedMethods, ", "))
	c.Header("Access-Control-Allow-Headers", strings.Join(hs.Config.Cors.AllowedHeaders, ", "))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", strconv.Itoa(hs.Config.Cors.MaxAge))

	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")

	// CSP adapted for APIs
	if c.Request.URL.Path != "/csp-report" { // Exclude endpoint from CSP
		csp := []string{
			"default-src 'none'", // API not serving web content
			"frame-ancestors 'none'",
			"base-uri 'none'",
			"form-action 'none'",
			"sandbox",
		}
		c.Header("Content-Security-Policy", strings.Join(csp, "; "))
	}

	// Force HTTPS in production
	if gin.Mode() == gin.ReleaseMode {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	c.Next()
}

func (hs *HeadersAndSecurityControllerImpl) CspReport(c *gin.Context) {
	var report map[string]interface{}
	if err := c.BindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report format"})
		return
	}

	// Log violation reports
	log.Printf("CSP Violation: %+v", report)
	c.Status(http.StatusOK)
}
