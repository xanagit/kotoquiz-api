package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type HealthController interface {
	HealthCheck(c *gin.Context)
}

type HealthControllerImpl struct {
	Service services.ApiHealthService
}

// Make sure that HealthControllerImpl implements HealthController
var _ HealthController = (*HealthControllerImpl)(nil)

func (hc *HealthControllerImpl) HealthCheck(c *gin.Context) {
	if err := hc.Service.Check(); err != nil {
		response := struct {
			Status string `json:"status"`
		}{
			Status: "DOWN",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "UP",
	}

	c.JSON(http.StatusOK, response)
}
