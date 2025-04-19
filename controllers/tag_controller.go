// Package controllers implements HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

// TagController defines the interface for tag-related HTTP endpoints
// It provides methods to manage categorization tags for words in the system
type TagController interface {
	// ListTags handles GET requests to retrieve all tags
	ListTags(c *gin.Context)
	// CreateTag handles POST requests to create a new tag
	CreateTag(c *gin.Context)
	// ReadTag handles GET requests to retrieve a specific tag by ID
	ReadTag(c *gin.Context)
	// UpdateTag handles PUT requests to update an existing tag
	UpdateTag(c *gin.Context)
	// DeleteTag handles DELETE requests to remove a tag
	DeleteTag(c *gin.Context)
}

// TagControllerImpl implements the TagController interface
// It depends on the LabelService for business logic operations on tags
type TagControllerImpl struct {
	Service services.LabelService
}

// Make sure that TagControllerImpl implements TagController
var _ TagController = (*TagControllerImpl)(nil)

// ListTags handles GET requests to retrieve all tags
// It returns a list of all tags in the system
//
// Responses:
//   - 200 OK with an array of tags on success
//   - 500 Internal Server Error if a server error occurs
func (tc *TagControllerImpl) ListTags(c *gin.Context) {
	labels, err := tc.Service.ListLabels(models.Tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

// CreateTag handles POST requests to create a new tag
// The tag data is expected in the request body as JSON
//
// Responses:
//   - 201 Created with the created tag on success
//   - 400 Bad Request if the tag data is invalid
//   - 500 Internal Server Error if a server error occurs
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

// ReadTag handles GET requests to retrieve a specific tag by ID
// The tag ID is expected as a URL parameter
//
// Responses:
//   - 200 OK with the tag data on success
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no tag with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (tc *TagControllerImpl) ReadTag(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	tag, err := tc.Service.ReadLabel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// UpdateTag handles PUT requests to update an existing tag
// The tag ID is expected as a URL parameter, and the updated tag data in the request body
//
// Responses:
//   - 200 OK with the updated tag on success
//   - 400 Bad Request if the ID or tag data is invalid
//   - 404 Not Found if no tag with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (tc *TagControllerImpl) UpdateTag(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	var tag models.Label
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag.ID = id
	if err := tc.Service.UpdateLabel(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// DeleteTag handles DELETE requests to remove a tag by ID
// The tag ID is expected as a URL parameter
//
// Responses:
//   - 204 No Content on successful deletion
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no tag with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (tc *TagControllerImpl) DeleteTag(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := tc.Service.DeleteLabel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
