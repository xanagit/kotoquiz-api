// Package controllers implements HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

// LevelController defines the interface for level-related HTTP endpoints
// It provides methods to manage proficiency and difficulty levels in the system
type LevelController interface {
	// ListLevels handles GET requests to retrieve all levels
	ListLevels(c *gin.Context)
	// CreateLevel handles POST requests to create a new level
	CreateLevel(c *gin.Context)
	// ReadLevel handles GET requests to retrieve a specific level by ID
	ReadLevel(c *gin.Context)
	// UpdateLevel handles PUT requests to update an existing level
	UpdateLevel(c *gin.Context)
	// DeleteLevel handles DELETE requests to remove a level
	DeleteLevel(c *gin.Context)
}

// LevelControllerImpl implements the LevelController interface
// It depends on the LevelService for business logic operations
type LevelControllerImpl struct {
	Service services.LevelService
}

// Make sure that LevelControllerImpl implements LevelController
var _ LevelController = (*LevelControllerImpl)(nil)

// ListLevels handles GET requests to retrieve all levels
// It returns a list of all level categories and their associated level names
//
// Responses:
//   - 200 OK with an array of levels on success
//   - 500 Internal Server Error if a server error occurs
func (lc *LevelControllerImpl) ListLevels(c *gin.Context) {
	levels, err := lc.Service.ListLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, levels)
}

// CreateLevel handles POST requests to create a new level
// The level data is expected in the request body as JSON
//
// Responses:
//   - 201 Created with the created level on success
//   - 400 Bad Request if the level data is invalid
//   - 500 Internal Server Error if a server error occurs
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

// ReadLevel handles GET requests to retrieve a specific level by ID
// The level ID is expected as a URL parameter
//
// Responses:
//   - 200 OK with the level data on success
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no level with the given ID exists
//   - 500 Internal Server Error if a server error occurs
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

// UpdateLevel handles PUT requests to update an existing level
// The level ID is expected as a URL parameter, and the updated level data in the request body
//
// Responses:
//   - 200 OK with the updated level on success
//   - 400 Bad Request if the ID or level data is invalid
//   - 404 Not Found if no level with the given ID exists
//   - 500 Internal Server Error if a server error occurs
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

// DeleteLevel handles DELETE requests to remove a level by ID
// The level ID is expected as a URL parameter
//
// Responses:
//   - 204 No Content on successful deletion
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no level with the given ID exists
//   - 500 Internal Server Error if a server error occurs
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
