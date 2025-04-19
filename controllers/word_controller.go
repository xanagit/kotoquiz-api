// Package controllers implements the HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type WordController interface {
	ReadWord(c *gin.Context)
	CreateWord(c *gin.Context)
	UpdateWord(c *gin.Context)
	DeleteWord(c *gin.Context)
}

// WordControllerImpl implements the WordController interface
// It depends on the WordService for business logic operations
type WordControllerImpl struct {
	Service services.WordService
}

// Make sure that WordControllerImpl implements WordController
var _ WordController = (*WordControllerImpl)(nil)

// ReadWord handles GET requests to retrieve a word by ID
// The word ID is expected as a URL parameter
//
// Responses:
//   - 200 OK with the word data on success
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no word with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (s *WordControllerImpl) ReadWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	word, err := s.Service.ReadWord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

// CreateWord handles POST requests to create a new word
// The word data is expected in the request body as JSON
//
// Responses:
//   - 201 Created with the created word on success
//   - 400 Bad Request if the word data is invalid
//   - 500 Internal Server Error if a server error occurs
func (s *WordControllerImpl) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Service.CreateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, word)
}

// UpdateWord handles PUT requests to update an existing word
// The word ID is expected as a URL parameter, and the updated word data in the request body
//
// Responses:
//   - 200 OK with the updated word on success
//   - 400 Bad Request if the ID or word data is invalid
//   - 404 Not Found if no word with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (s *WordControllerImpl) UpdateWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}

	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID in the URL matches the ID in the body
	word.ID = id
	if err := s.Service.UpdateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

// DeleteWord handles DELETE requests to remove a word by ID
// The word ID is expected as a URL parameter
//
// Responses:
//   - 204 No Content on successful deletion
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no word with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (s *WordControllerImpl) DeleteWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := s.Service.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
