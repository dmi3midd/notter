package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStore    domain.UserRepository
	tokenService domain.TokenService
}

func NewUserService(userStore domain.UserRepository, tokenService domain.TokenService) *UserService {
	return &UserService{
		userStore:    userStore,
		tokenService: tokenService,
	}
}

func (us *UserService) Registration(ctx context.Context, username, email, password string) (*domain.UserData, error) {
	op := "user.service-Registration"
	candidate, err := us.userStore.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if candidate != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExist)
	}

	id := uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user, err := us.userStore.Create(ctx, id, username, email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	userDto := domain.NewUserDto(user)
	tokens, err := us.tokenService.GenerateTokens(*userDto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := us.tokenService.SaveToken(ctx, userDto.Id, tokens.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (us *UserService) Login(ctx context.Context, email, password string) (*domain.UserData, error) {
	op := "user.service-Login"
	user, err := us.userStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidPw)
	}

	userDto := domain.NewUserDto(user)
	tokens, err := us.tokenService.GenerateTokens(*userDto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := us.tokenService.SaveToken(ctx, userDto.Id, tokens.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (us *UserService) Logout(ctx context.Context, refreshToken string) error {
	op := "user,service-Logout"
	err := us.tokenService.RemoveToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (us *UserService) Refresh(ctx context.Context, refreshToken string) (*domain.UserData, error) {
	op := "user.service-Refresh"
	userFromToken := us.tokenService.ValidateRefreshToken(refreshToken)
	_, err := us.tokenService.FindToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrUnuthorized)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if userFromToken == nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUnuthorized)
	}

	user, err := us.userStore.GetByEmail(ctx, userFromToken.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	userDto := domain.NewUserDto(user)
	tokens, err := us.tokenService.GenerateTokens(*userDto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = us.tokenService.SaveToken(ctx, userDto.Id, tokens.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil

}
