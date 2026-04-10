package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	store domain.TokenRepository
	cfg   config.JWTConfig
}

func NewTokenService(cfg config.JWTConfig, store domain.TokenRepository) *TokenService {
	return &TokenService{
		store: store,
		cfg:   cfg,
	}
}

type UserClaims struct {
	domain.UserDto
	jwt.RegisteredClaims
}

func (s *TokenService) GenerateTokens(payload domain.UserDto) (*domain.TokenPair, error) {
	op := "token.service-GenerateToken"
	accessSecret := []byte(s.cfg.AccessSecret)
	refreshSecret := []byte(s.cfg.RefreshSecret)
	accessExpiry := s.cfg.AccessExpiry
	refreshExpiry := s.cfg.RefreshExpiry

	accessClaims := UserClaims{
		UserDto: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := UserClaims{
		UserDto: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *TokenService) SaveToken(ctx context.Context, userId, refreshToken string) error {
	op := "token.service-SaveToken"
	if err := s.store.Update(ctx, userId, refreshToken); err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			if err := s.store.Create(ctx, userId, refreshToken); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *TokenService) ValidateAccessToken(accessToken string) *domain.UserDto {
	op := "token.service-validateAccessToken"
	accessSecret := []byte(s.cfg.AccessSecret)

	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w, %v", op, domain.ErrSignMethod, token.Header["alg"])
		}
		return accessSecret, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		return &claims.UserDto
	}

	return nil
}

func (s *TokenService) ValidateRefreshToken(refreshToken string) *domain.UserDto {
	op := "token.service-validateRefreshToken"
	refreshSecret := []byte(s.cfg.RefreshSecret)

	token, err := jwt.ParseWithClaims(refreshToken, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w, %v", op, domain.ErrSignMethod, token.Header["alg"])
		}
		return refreshSecret, nil
	})

	if err != nil {
		fmt.Printf("Refresh token validation error: %v\n", err)
		return nil
	}

	if !token.Valid {
		fmt.Println("Refresh token is invalid")
		return nil
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		return &claims.UserDto
	}

	return nil
}

func (s *TokenService) RemoveToken(ctx context.Context, refreshToken string) error {
	op := "token.service-RemoveToken"
	err := s.store.Delete(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *TokenService) FindToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	op := "token.service-FindToken"
	token, err := s.store.GetToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
