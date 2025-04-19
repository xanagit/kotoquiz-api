// Package controllers implements HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

// HealthController defines the interface for API health check endpoints
// This controller is responsible for providing application health status information
type HealthController interface {
	// HealthCheck handles GET requests to check the API's health status
	HealthCheck(c *gin.Context)
}

// HealthControllerImpl implements the HealthController interface
// It depends on the ApiHealthService for retrieving health status information
type HealthControllerImpl struct {
	Service services.ApiHealthService
}

// Make sure that HealthControllerImpl implements HealthController
var _ HealthController = (*HealthControllerImpl)(nil)

// HealthCheck handles GET requests to the /health endpoint
// It returns basic health information about the API, including database connectivity
//
// Responses:
//   - 200 OK with health status details on success
//   - 503 Service Unavailable if any critical component is unhealthy
func (hc *HealthControllerImpl) HealthCheck(c *gin.Context) {
	if err := hc.Service.Check(); err != nil {
		response := struct {
			Status string `json:"status"`
		}{
			Status: "DOWN",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "UP",
	}

	c.JSON(http.StatusOK, response)
}
