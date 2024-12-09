package middlewares

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

type KotoquizRole string

var (
	UserRole  KotoquizRole = "user"
	AdminRole KotoquizRole = "admin"
)

type AuthMiddleware interface {
	AuthRequired() gin.HandlerFunc
	RequireRoles(roles ...string) gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
	config       *config.KeycloakConfig
}

// Claims structure for JWT token claims
type Claims struct {
	Subject           string `json:"sub"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// Make sure that AuthMiddlewareImpl implements AuthMiddleware
var _ AuthMiddleware = (*AuthMiddlewareImpl)(nil)

func NewAuthMiddleware(cfg *config.KeycloakConfig) (*AuthMiddlewareImpl, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &AuthMiddlewareImpl{
		provider:     provider,
		verifier:     verifier,
		oauth2Config: oauth2Config,
		config:       cfg,
	}, nil
}

// AuthRequired is the middleware that checks for valid JWT token
func (am *AuthMiddlewareImpl) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := bearerToken[1]
		idToken, err := am.verifier.Verify(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})
			return
		}

		// Store claims in context for later use
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRoles middleware checks if the user has the required roles
func (am *AuthMiddlewareImpl) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no claims found"})
			return
		}

		claims := claimsInterface.(Claims)
		userRoles := claims.RealmAccess.Roles

		log.Print("User roles:")
		log.Print(userRoles)

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if requiredRole == userRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}

// GetUserIDFromToken extracts user ID from Keycloak token
func GetUserIDFromToken(c *gin.Context) (string, error) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return "", fmt.Errorf("no claims found")
	}

	claims := claimsInterface.(Claims)
	return claims.Subject, nil // "sub" claim contains keycloak user ID
}
