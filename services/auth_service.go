package services

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"time"
)

type AuthService interface {
	GenerateToken(userID uuid.UUID, role models.UserRole, tokenType models.TokenType) (string, error)
	ValidateToken(tokenString string) (*models.TokenClaims, error)
	RefreshTokens(refreshToken string) (*models.TokenPair, error)
}

type AuthServiceImpl struct {
	JwtSecret []byte
	// Injection of UserService to check credentials
	UserService UserService
}

func (s *AuthServiceImpl) GenerateToken(userID uuid.UUID, role models.UserRole, tokenType models.TokenType) (string, error) {
	// Determine the validity period based on the token type
	var ttl time.Duration
	if tokenType == models.AccessTokenType {
		ttl = 30 * time.Minute // Access token valid for 30 minutes
	} else {
		ttl = 720 * time.Hour // Refresh token valid for 30 days
	}

	claims := models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "kotoquiz-api",
			Subject:   userID.String(),
		},
		UserID: userID,
		Role:   role,
		Type:   tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JwtSecret)
}

func (s *AuthServiceImpl) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthServiceImpl) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	// Validate the refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if it is a refresh token
	if claims.Type != models.RefreshTokenType {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	// Check if the user still exists and get their information
	user, err := s.UserService.ReadUser(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate a new access token
	accessToken, err := s.GenerateToken(user.ID, user.Role, models.AccessTokenType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Générer un nouveau refresh token
	newRefreshToken, err := s.GenerateToken(user.ID, user.Role, models.RefreshTokenType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
