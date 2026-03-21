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
	// If necessary, return ErrBoardNotFound
	GetBoard(
		ctx context.Context,
		boardId string,
	) (*Board, error)

	GetBoardsByUserId(
		ctx context.Context,
		userId string,
	) ([]Board, error)

	CreateBoard(
		ctx context.Context,
		board *Board,
	) error

	// If necessary, return ErrBoardNotFound
	UpdateBoard(
		ctx context.Context,
		board *Board,
	) error

	DeleteBoard(
		ctx context.Context,
		boardId string,
	) error
}

type BoardService interface {
	GetBoards(
		ctx context.Context,
		userId string,
	) ([]BoardDto, error)

	CreateBoard(
		ctx context.Context,
		userId string,
		title string,
	) error

	UpdateBoard(
		ctx context.Context,
		boardId string,
		title string,
	) error

	DeleteBoard(
		ctx context.Context,
		boardId string,
	) error
}
