// Package controllers implements HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type RegistrationController interface {
	RegisterUser(c *gin.Context)
}

type RegistrationControllerImpl struct {
	Service services.RegistrationService
}

// RegistrationRequest represents the data needed to register a new user
// It contains all required fields for user creation in the authentication system
type RegistrationRequest struct {
	// Username is the unique identifier for the user's account
	Username string `json:"username" binding:"required"`
	// Email is the user's email address for account recovery and notifications
	Email string `json:"email" binding:"required,email"`
	// Password is the user's account password (will be hashed before storage)
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterUser handles POST requests to register new users
// It validates the request data and forwards it to the registration service
//
// Responses:
//   - 201 Created on successful user registration
//   - 400 Bad Request if the request data is invalid
//   - 409 Conflict if a user with the same username or email already exists
//   - 500 Internal Server Error if registration fails for other reasons
func (rc *RegistrationControllerImpl) RegisterUser(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rc.Service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
