package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type CategoryController interface {
	ListCategories(c *gin.Context)
	CreateCategory(c *gin.Context)
	ReadCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type CategoryControllerImpl struct {
	Service services.LabelService
}

func (cc *CategoryControllerImpl) ListCategories(c *gin.Context) {
	labels, err := cc.Service.ListLabels("CATEGORY")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (cc *CategoryControllerImpl) CreateCategory(c *gin.Context) {
	var category models.Label
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category.Type = "CATEGORY"
	if err := cc.Service.CreateLabel(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (cc *CategoryControllerImpl) ReadCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.Service.ReadLabel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *CategoryControllerImpl) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Label
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category.ID = FromStrToUuid(id)
	category.Type = "CATEGORY"
	if err := cc.Service.UpdateLabel(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *CategoryControllerImpl) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteLabel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
