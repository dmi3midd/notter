package repositories

import (
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

func (r *TokenRepository) GetToken(refreshToken string) (*domain.Token, error) {
	query := "SELECT * FROM tokens WHERE refresh_token = $1"
	var token domain.Token
	err := r.store.Get(&token, query, refreshToken)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) Create(userId, refreshToken string) (*domain.Token, error) {
	query := `INSERT INTO tokens (user_id, refresh_token)
			  VALUES ($1, $2) RETURNING *`
	var token domain.Token
	err := r.store.Get(&token, query, userId, refreshToken)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) Delete(refreshToken string) (*domain.Token, error) {
	query := "DELETE FROM tokens WHERE refresh_token = $1 RETURNING *"
	var token domain.Token
	err := r.store.Get(&token, query, refreshToken)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) Update(userId, refreshToken string) (*domain.Token, error) {
	query := "UPDATE tokens SET refresh_token = $1 WHERE user_id = $2 RETURNING *"
	var token domain.Token
	err := r.store.Get(&token, query, refreshToken, userId)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
