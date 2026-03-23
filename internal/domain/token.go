package domain

import (
	"context"
	"errors"
)

var ErrTokenNotFound error = errors.New("token not found")
var ErrSignMethod error = errors.New("unexpected signing")

type Token struct {
	UserId       string `json:"userId" db:"user_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenRepository interface {
	// GetToken retrieves a Token entity by its refresh token string.
	// It returns ErrTokenNotFound if no token are found.
	GetToken(ctx context.Context, refreshToken string) (*Token, error)
	// Create creates a Token entity.
	Create(ctx context.Context, userId, refreshToken string) error
	// Delete deletes a Token entity by its refresh token string.
	Delete(ctx context.Context, refreshToken string) error
	// Updated updates a Token entity.
	// It returns ErrTokenNotFound if no token are found.
	Update(ctx context.Context, userId, refreshToken string) error
}

type TokenService interface {
	// GenerateTokens generates pair with access and refresh tokens.
	GenerateTokens(payload UserDto) (*TokenPair, error)
	// SaveToken updates an existing refresh token for the user or creates a new one.
	SaveToken(ctx context.Context, userId, refreshToken string) error
	// ValidateAccessToken validates access token and returns a UserDto struct.
	// It returns nil if validation go wrong.
	ValidateAccessToken(accessToken string) *UserDto
	// ValidateRefreshToken validates refresh token and returns a UserDto struct.
	// It returns nil if validation go wrong.
	ValidateRefreshToken(refreshToken string) *UserDto
	// RemoveToken removes refresh token.
	// It returns nil if validation go wrong.
	RemoveToken(ctx context.Context, refreshToken string) error
	// FindToken finds and returns a Token entity by its refresh token string.
	// It returns ErrTokenNotFound if no token are found.
	FindToken(ctx context.Context, refreshToken string) (*Token, error)
}
