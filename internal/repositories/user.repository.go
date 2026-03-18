package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	op := "user.repository-GetById"
	query := "SELECT * FROM users WHERE id = $1"
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	op := "user.repository-GetByEmail"
	query := "SELECT * FROM users WHERE email = $1"
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, id, username, email, hashedPassword string) (*domain.User, error) {
	op := "user.repositroy-Create"
	query := `INSERT INTO users (id, username, email, hashed_password) 
			  VALUES ($1, $2, $3, $4) RETURNING *`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id, username, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}
