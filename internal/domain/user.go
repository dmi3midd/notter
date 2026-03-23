package domain

import (
	"context"
	"errors"
	"time"
)

var ErrUserNotFound error = errors.New("user not found")
var ErrUserAlreadyExist error = errors.New("user already exist")
var ErrUnauthorized error = errors.New("user unauthorized")
var ErrInvalidPw error = errors.New("invalid password")

type User struct {
	Id             string    `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	HashedPassword string    `json:"hashedPassword" db:"hashed_password"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type UserData struct {
	User         UserDto `json:"user"`
	RefreshToken string  `json:"refreshToken"`
	AccessToken  string  `json:"accessToken"`
}

type UserDto struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserDto(user *User) *UserDto {
	return &UserDto{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}

type UserRepository interface {
	// GetById retrieves a User entity by its id.
	// It returns ErrUserNotFound if no user are found.
	GetById(ctx context.Context, id string) (*User, error)
	// GetByEmail retrieves a User entity by its email.
	// It returns ErrUserNotFound if no user are found.
	GetByEmail(ctx context.Context, email string) (*User, error)
	// Create creates a User entity and returns it.
	Create(ctx context.Context, id, username, email, hashedPassword string) (*User, error)
}

type UserService interface {
	// Registration performs user registration and returns UserData struct.
	// It returns ErrUserAlreadyExist if the user exist.
	Registration(ctx context.Context, username, email, password string) (*UserData, error)
	// Login performs user login and returns UserData struct.
	// It returns ErrUserNotFound if no user are found.
	// It returns ErrInvalidPw if the password is invalid.
	Login(ctx context.Context, email, password string) (*UserData, error)
	// Logout performs logout user.
	Logout(ctx context.Context, refreshToken string) error
	// Refresh performs refresh access and refresh tokens and returns UserData struct.
	// It returns ErrUnuthorized if no refresh token are found or ValidateRefreshToken returned nil
	Refresh(ctx context.Context, refreshToken string) (*UserData, error)
}
