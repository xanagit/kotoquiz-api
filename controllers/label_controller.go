package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type LabelController interface {
	CreateLabel(c *gin.Context)
	ReadLabel(c *gin.Context)
	UpdateLabel(c *gin.Context)
	DeleteLabel(c *gin.Context)
}

type LabelControllerImpl struct {
	Service services.LabelService
}

func (tc *LabelControllerImpl) CreateLabel(c *gin.Context) {
	var tag models.Label
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.Service.CreateLabel(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tag)
}

func (tc *LabelControllerImpl) ReadLabel(c *gin.Context) {
	id := c.Param("id")
	tag, err := tc.Service.ReadLabel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (tc *LabelControllerImpl) UpdateLabel(c *gin.Context) {
	id := c.Param("id")
	var tag models.Label
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag.ID = FromStrToUuid(id)
	tag.Type = "TAG"
	if err := tc.Service.UpdateLabel(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (tc *LabelControllerImpl) DeleteLabel(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteLabel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
