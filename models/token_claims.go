package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAppUser UserRole = "APP_USER"
	RoleTech    UserRole = "TECH"
)

type TokenType string

const (
	AccessTokenType  TokenType = "ACCESS"
	RefreshTokenType TokenType = "REFRESH"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userId"`
	Role   UserRole  `json:"role"`
	Type   TokenType `json:"type"` // Type of token: ACCESS or REFRESH
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
