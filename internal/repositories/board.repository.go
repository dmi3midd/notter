package repositories

import (
	"context"
	"database/sql"
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
	query := "SELECT * FROM boards WHERE id = $1"
	var board domain.Board
	if err := r.db.GetContext(ctx, &board, query, boardId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	query := `INSERT INTO boards (id, user_id, title, notes, created_at, updated_at) 
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
	query := `UPDATE boards SET title = :title, updated_at = :updated_at WHERE id = :id`
	res, err := r.db.NamedExecContext(ctx, query, board)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrBoardNotFound)
	}
	return nil
}

func (r *BoardRepository) DeleteBoard(
	ctx context.Context,
	boardId string,
) error {
	const op = "board.repository-DeleteBoard"

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	queryToDeleteBoard := "DELETE FROM boards WHERE id = $1"
	res, err := tx.ExecContext(ctx, queryToDeleteBoard, boardId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrBoardNotFound)
	}

	queryToDeleteNotes := "DELETE FROM notes WHERE board_id = $1"
	if _, err := tx.ExecContext(ctx, queryToDeleteNotes, boardId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
