package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type LevelNameController interface {
	ListLevelNames(c *gin.Context)
	CreateLevelName(c *gin.Context)
	ReadLevelName(c *gin.Context)
	UpdateLevelName(c *gin.Context)
	DeleteLevelName(c *gin.Context)
}

type LevelNameControllerImpl struct {
	Service services.LabelService
}

func (lnc *LevelNameControllerImpl) ListLevelNames(c *gin.Context) {
	cid := c.Param("id")
	labels, err := lnc.Service.ListLabelsOfCategory("LEVEL_NAME", cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (lnc *LevelNameControllerImpl) CreateLevelName(c *gin.Context) {
	cid := c.Param("id")
	var levelName models.Label
	if err := c.ShouldBindJSON(&levelName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	levelName.Type = "LEVEL_NAME"
	if err := lnc.Service.CreateCategoryLabel(&levelName, cid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, levelName)
}

func (lnc *LevelNameControllerImpl) ReadLevelName(c *gin.Context) {
	id := c.Param("id")
	levelName, err := lnc.Service.ReadLabel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levelName)
}

func (lnc *LevelNameControllerImpl) UpdateLevelName(c *gin.Context) {
	id := c.Param("id")
	var levelName models.Label
	if err := c.ShouldBindJSON(&levelName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	levelName.ID = FromStrToUuid(id)
	levelName.Type = "LEVEL_NAME"
	if err := lnc.Service.UpdateLabel(&levelName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levelName)
}

func (lnc *LevelNameControllerImpl) DeleteLevelName(c *gin.Context) {
	id := c.Param("id")
	if err := lnc.Service.DeleteLabel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
