package domain

import (
	"context"
	"errors"
)

var ErrTokenNotFound error = errors.New("token not found")
var ErrSignMethod error = errors.New("unexpected signing method")

type Token struct {
	UserId       string `json:"userId" db:"user_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenRepository interface {
	GetToken(ctx context.Context, refreshToken string) (*Token, error)
	Create(ctx context.Context, userId, refreshToken string) error
	Delete(ctx context.Context, refreshToken string) error
	Update(ctx context.Context, userId, refreshToken string) error
}

type TokenService interface {
	GenerateTokens(payload UserDto) (*TokenPair, error)
	SaveToken(ctx context.Context, userId, refreshToken string) error
	ValidateAccessToken(accessToken string) *UserDto
	ValidateRefreshToken(refreshToken string) *UserDto
	RemoveToken(ctx context.Context, refreshToken string) error
	FindToken(ctx context.Context, refreshToken string) (*Token, error)
}
