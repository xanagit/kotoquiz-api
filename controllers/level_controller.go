package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type LevelController interface {
	ListLevels(c *gin.Context)
	CreateLevel(c *gin.Context)
	ReadLevel(c *gin.Context)
	UpdateLevel(c *gin.Context)
	DeleteLevel(c *gin.Context)
}

type LevelControllerImpl struct {
	Service services.LevelService
}

// Make sure that LevelControllerImpl implements LevelController
var _ LevelController = (*LevelControllerImpl)(nil)

func (lc *LevelControllerImpl) ListLevels(c *gin.Context) {
	levels, err := lc.Service.ListLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levels)
}

func (lc *LevelControllerImpl) CreateLevel(c *gin.Context) {
	var level models.Level
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := lc.Service.CreateLevel(&level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, level)
}

func (lc *LevelControllerImpl) ReadLevel(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	level, err := lc.Service.ReadLevel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, level)
}

func (lc *LevelControllerImpl) UpdateLevel(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	var level models.Level
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	level.ID = id

	if err := lc.Service.UpdateLevel(&level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, level)
}

func (lc *LevelControllerImpl) DeleteLevel(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := lc.Service.DeleteLevel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
