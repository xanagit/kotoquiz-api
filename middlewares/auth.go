package middlewares

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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
	logger       *zap.Logger
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

func NewAuthMiddleware(cfg *config.KeycloakConfig, log *zap.Logger) (*AuthMiddlewareImpl, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		log.Error("Failed to create OIDC provider",
			zap.String("issuerURL", cfg.IssuerURL),
			zap.Error(err))
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

	log.Info("Auth middleware initialized successfully",
		zap.String("clientID", cfg.ClientID),
		zap.String("issuerURL", cfg.IssuerURL))

	return &AuthMiddlewareImpl{
		provider:     provider,
		verifier:     verifier,
		oauth2Config: oauth2Config,
		config:       cfg,
		logger:       log,
	}, nil
}

// AuthRequired is the middleware that checks for valid JWT token
func (am *AuthMiddlewareImpl) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			am.logger.Warn("Authorization header missing in request",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("clientIP", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			am.logger.Warn("Invalid authorization header format",
				zap.String("path", c.Request.URL.Path),
				zap.String("header", authHeader),
				zap.String("clientIP", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := bearerToken[1]
		idToken, err := am.verifier.Verify(c.Request.Context(), token)
		if err != nil {
			am.logger.Warn("Invalid token verification",
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
				zap.String("clientIP", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			am.logger.Error("Failed to parse token claims",
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
				zap.String("clientIP", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})
			return
		}

		// Store claims in context for later use
		c.Set("claims", claims)

		am.logger.Debug("Successful authentication",
			zap.String("subject", claims.Subject),
			zap.String("username", claims.PreferredUsername),
			zap.String("path", c.Request.URL.Path))

		c.Next()
	}
}

// RequireRoles middleware checks if the user has the required roles
func (am *AuthMiddlewareImpl) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			am.logger.Warn("No claims found in context",
				zap.String("path", c.Request.URL.Path),
				zap.String("clientIP", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no claims found"})
			return
		}

		claims := claimsInterface.(Claims)
		userRoles := claims.RealmAccess.Roles

		am.logger.Debug("Checking user roles",
			zap.Strings("userRoles", userRoles),
			zap.Strings("requiredRoles", roles),
			zap.String("username", claims.PreferredUsername),
			zap.String("path", c.Request.URL.Path))

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
			am.logger.Warn("Insufficient permissions",
				zap.Strings("userRoles", userRoles),
				zap.Strings("requiredRoles", roles),
				zap.String("username", claims.PreferredUsername),
				zap.String("path", c.Request.URL.Path))
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
