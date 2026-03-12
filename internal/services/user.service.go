package services

import (
	"database/sql"
	"errors"
	"log"

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
	if err := us.tokenService.SaveToken(userDto.Id, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (us *UserService) Login(email, password string) (*domain.UserData, error) {
	user, err1 := us.store.GetByEmail(email)
	if err1 != nil {
		return nil, err1
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	userDto := domain.NewUserDto(user)
	tokens, err := us.tokenService.GenerateTokens(*userDto)
	if err != nil {
		return nil, err
	}
	if err := us.tokenService.SaveToken(userDto.Id, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (us *UserService) Logout(refreshToken string) error {
	err := us.tokenService.RemoveToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) Refresh(refreshToken string) (*domain.UserData, error) {
	userFromToken := us.tokenService.ValidateRefreshToken(refreshToken)
	foundToekn, err1 := us.tokenService.FindToken(refreshToken)
	if err1 != nil {
		return nil, err1
	}
	if userFromToken == nil || foundToekn == nil {
		log.Println(userFromToken)
		log.Println(foundToekn)
		return nil, errors.New("Unuthorizaed")
	}

	user, err2 := us.store.GetByEmail(userFromToken.Email)
	if err2 != nil {
		return nil, err2
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	userDto := domain.NewUserDto(user)
	tokens, err3 := us.tokenService.GenerateTokens(*userDto)
	if err3 != nil {
		return nil, err3
	}
	us.tokenService.SaveToken(userDto.Id, tokens.RefreshToken)
	return &domain.UserData{
		User:         *userDto,
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil

}
