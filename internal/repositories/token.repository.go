package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	store *sqlx.DB
}

func NewTokenRepo(store *sqlx.DB) *TokenRepository {
	return &TokenRepository{
		store: store,
	}
}

func (r *TokenRepository) GetToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	op := "token.repository-GetToken"
	query := "SELECT * FROM tokens WHERE refresh_token = $1"
	var token domain.Token
	err := r.store.GetContext(ctx, &token, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &token, nil
}

func (r *TokenRepository) Create(ctx context.Context, userId, refreshToken string) error {
	op := "token.repository-Create"
	query := `INSERT INTO tokens (user_id, refresh_token)
			  VALUES ($1, $2) RETURNING *`
	var token domain.Token
	err := r.store.GetContext(ctx, &token, query, userId, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TokenRepository) Delete(ctx context.Context, refreshToken string) error {
	op := "token.repository-Delete"
	query := "DELETE FROM tokens WHERE refresh_token = $1 RETURNING *"
	var token domain.Token
	err := r.store.GetContext(ctx, &token, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TokenRepository) Update(ctx context.Context, userId, refreshToken string) error {
	op := "token.repository-Update"
	query := "UPDATE tokens SET refresh_token = $1 WHERE user_id = $2 RETURNING *"
	var token domain.Token
	err := r.store.GetContext(ctx, &token, query, refreshToken, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
