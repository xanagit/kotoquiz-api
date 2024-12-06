package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthController interface {
	Login(c *gin.Context)
}

type AuthControllerImpl struct {
	AuthService services.AuthService
	UserService services.UserService
}

func (as *AuthControllerImpl) Login(c *gin.Context) {
	// Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check credentials and get the role
	user, err := as.UserService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	//Generate token with appropriate role
	token, err := as.AuthService.GenerateToken(user.ID, user.Role, models.AccessTokenType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// controllers/auth_controller.go

func (ac *AuthControllerImpl) RefreshToken(c *gin.Context) {
	// Get the refresh token from the Authorization header
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token provided"})
		return
	}

	// Remove "Bearer " prefix if present
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	// Call the service to generate a new pair of tokens
	newTokens, err := ac.AuthService.RefreshTokens(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, newTokens)
}
