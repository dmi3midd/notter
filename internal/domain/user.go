package domain

import (
	"context"
	"errors"
	"time"
)

var ErrUserNotFound error = errors.New("user not found")
var ErrUserAlreadyExist error = errors.New("user already exist")
var ErrUnuthorized error = errors.New("user unauthorized")
var ErrInvalidPw error = errors.New("invalid password")

type User struct {
	Id             string    `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	HashedPassword string    `json:"hashed_password" db:"hashed_password"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type UserData struct {
	User         UserDto `json:"user"`
	RefreshToken string  `json:"refresh_token"`
	AccessToken  string  `json:"access_token"`
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
	GetById(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, id, username, email, hashedPassword string) (*User, error)
}
