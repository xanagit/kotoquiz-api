package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type TagController interface {
	ListTags(c *gin.Context)
	CreateTag(c *gin.Context)
	ReadTag(c *gin.Context)
	UpdateTag(c *gin.Context)
	DeleteTag(c *gin.Context)
}

type TagControllerImpl struct {
	Service services.LabelService
}

func (tc *TagControllerImpl) ListTags(c *gin.Context) {
	labels, err := tc.Service.ListLabels("TAG")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (tc *TagControllerImpl) CreateTag(c *gin.Context) {
	var tag models.Label
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.Service.CreateLabel(&tag, models.Tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tag)
}

func (tc *TagControllerImpl) ReadTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := tc.Service.ReadLabel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (tc *TagControllerImpl) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Label
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag.ID = FromStrToUuid(id)
	if err := tc.Service.UpdateLabel(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (tc *TagControllerImpl) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteLabel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
