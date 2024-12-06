package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
	"strings"
)

type AuthenticationMiddleware interface {
	AuthRequired(requiredRole models.UserRole) gin.HandlerFunc
}

type AuthenticationMiddlewareImpl struct {
	AuthService services.AuthService
}

func (am *AuthenticationMiddlewareImpl) AuthRequired(requiredRole models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := am.AuthService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userId", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
