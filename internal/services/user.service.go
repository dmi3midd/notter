package services

import (
	"database/sql"
	"errors"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store        domain.UserRepository
	tokenService *TokenService
}

func NewUserService(store domain.UserRepository, tokenService *TokenService) *UserService {
	return &UserService{
		store:        store,
		tokenService: tokenService,
	}
}

func (us *UserService) Registration(username, email, password string) (*domain.UserData, error) {
	candidate, err1 := us.store.GetByEmail(email)
	if err1 != nil && err1 != sql.ErrNoRows {
		return nil, err1
	}

	if candidate != nil {
		return nil, errors.New("User with this email already exists")
	}

	id := uuid.NewString()
	hashedPassword, err2 := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err2 != nil {
		return nil, err2
	}
	user, err3 := us.store.Create(id, username, email, string(hashedPassword))
	if err3 != nil {
		return nil, err3
	}
	userDto := domain.NewUserDto(user)
	tokens, err4 := us.tokenService.GenerateTokens(*userDto)
	if err4 != nil {
		return nil, err4
	}
	_, err5 := us.tokenService.SaveToken(userDto.Id, tokens.RefreshToken)
	if err5 != nil {
		return nil, err5
	}

	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}
