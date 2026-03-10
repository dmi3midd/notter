package domain

import (
	"time"
)

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
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(id, username, email, hashedPassword string) (*User, error)
}
