package repositories

import (
	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	store *sqlx.DB
}

func NewUserRepo(store *sqlx.DB) *UserRepository {
	return &UserRepository{
		store: store,
	}
}

func (r *UserRepository) GetById(id string) (*domain.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	var user domain.User
	err := r.store.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	var user domain.User
	err := r.store.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(id, username, email, hashedPassword string) (*domain.User, error) {
	query := `INSERT INTO users (id, username, email, hashed_password) 
			  VALUES ($1, $2, $3, $4) RETURNING *`
	var user domain.User
	err := r.store.Get(&user, query, id, username, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
