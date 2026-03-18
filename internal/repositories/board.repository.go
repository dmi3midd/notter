package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BoardRepository struct {
	db *sqlx.DB
}

func NewBoardRepo(db *sqlx.DB) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

func (r *BoardRepository) GetBoard(
	ctx context.Context,
	boardId string,
) (*domain.Board, error) {
	op := "board.repository-GetBoard"
	query := "SELEct * FROM boards WHERE board_id = $1"
	var board domain.Board
	if err := r.db.GetContext(ctx, &board, query, boardId); err != nil {
		if errors.Is(err, domain.ErrBoardNotFound) {
			return nil, domain.ErrBoardNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &board, nil
}

func (r *BoardRepository) GetBoardsByUserId(
	ctx context.Context,
	userId string,
) ([]domain.Board, error) {
	op := "board.repository-GetBoardsByUserId"
	query := "SELECT * FROM boards WHERE user_id = $1"
	var boards []domain.Board
	if err := r.db.SelectContext(ctx, &boards, query, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return boards, nil
}

func (r *BoardRepository) CreateBoard(
	ctx context.Context,
	board *domain.Board,
) error {
	op := "board.repository-CreateBoard"
	query := `INSERT INTO boards 
			         (id, user_id, title, notes, created_at, updated_at)
			  VALUES (:id, :user_id, :title, :notes, :created_at, :updated_at)
	`
	if _, err := r.db.NamedExecContext(ctx, query, board); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *BoardRepository) UpdateBoard(
	ctx context.Context,
	board *domain.Board,
) error {
	op := "board.repository-UpdateBoard"
	query := `UPDATE boards SET
			  	title = :title, notes = :notes, updated_at = :updated_at
			  WHERE id = :id
	`
	if _, err := r.db.NamedExecContext(ctx, query, board); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *BoardRepository) DeleteBoard(
	ctx context.Context,
	boardId string,
) error {
	op := "board.repository-DeleteBoard"
	query := "DELETE FROM boards WHERE id = $1"
	if _, err := r.db.ExecContext(ctx, query, boardId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
