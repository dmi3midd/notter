package domain

import (
	"context"
	"errors"
	"time"
)

var ErrBoardNotFound error = errors.New("board not found")

type Board struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Notes     int       `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type BoardDto struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Notes     int       `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToBoardDto(board *Board) *BoardDto {
	return &BoardDto{
		Id:        board.Id,
		Title:     board.Title,
		Notes:     board.Notes,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
	}
}

type BoardRepository interface {
	// GetBoard retrieves a single board by its id.
	// It returns domain.ErrBoardNotFound if no board are found.
	GetBoard(ctx context.Context, boardId string) (*Board, error)
	// GetBoardsByUserId retieves a user's boards.
	// It returns an empty slice if no boards are found.
	GetBoardsByUserId(ctx context.Context, userId string) ([]Board, error)
	// CreateBoard inserts a new board.
	CreateBoard(ctx context.Context, board *Board) error
	// UpdateBoard updates the existing board.
	// It returns domain.ErrBoardNotFound if no board are found.
	UpdateBoard(ctx context.Context, board *Board) error
	// DeleteBoard removes a board and all its associated notes in a single transaction.
	// It returns domain.ErrBoardNotFound if no board are found.
	DeleteBoard(ctx context.Context, boardId string) error
}

type BoardService interface {
	// GetBoards finds a user's boards.
	// It returns ErrUserNotFound if no user are found.
	// It returns an empty slice if no boards are found.
	GetBoards(ctx context.Context, userId string) ([]BoardDto, error)
	// CreateBoard creates a new board.
	// It returns ErrUserNotFound if no user are found.
	CreateBoard(ctx context.Context, userId string, title string) error
	// UpdateBoard modifies the existing board.
	// It returns domain.ErrBoardNotFound if no board are found.
	UpdateBoard(ctx context.Context, boardId string, title string) error
	// DeleteBoard removes the board.
	DeleteBoard(ctx context.Context, boardId string) error
}
