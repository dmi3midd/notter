package services

import (
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

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (ts *TokenService) GenerateTokens(payload domain.UserDto) (*TokenPair, error) {
	accessSecret := []byte(ts.cfg.JWT_ACCESS_SECRET)
	refreshSecret := []byte(ts.cfg.JWT_REFRESH_SECRET)
	accessExpiry := ts.cfg.AccessExpiry
	refreshExpiry := ts.cfg.RefreshExpiry

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
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *TokenService) SaveToken(userId, refreshToken string) (*domain.Token, error) {
	_, err1 := ts.store.GetToken(refreshToken)
	if err1 != nil {
		token, err2 := ts.store.Create(userId, refreshToken)
		if err2 != nil {
			return nil, err2
		}
		return token, nil
	}
	token, err3 := ts.store.Update(userId, refreshToken)
	if err3 != nil {
		return nil, err3
	}
	return token, nil
}

func (ts *TokenService) ValidateAccessToken(accessToken string) *domain.UserDto {
	accessSecret := []byte(ts.cfg.JWT_ACCESS_SECRET)

	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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

func (ts *TokenService) ValidateRefreshToken(refreshToken string) *domain.UserDto {
	refreshSecret := []byte(ts.cfg.JWT_REFRESH_SECRET)

	token, err := jwt.ParseWithClaims(refreshToken, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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

func (ts *TokenService) RemoveToken(refreshToken string) (*domain.Token, error) {
	token, err := ts.store.Delete(refreshToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (ts *TokenService) FindToken(refreshToken string) (*domain.Token, error) {
	token, err := ts.store.GetToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
