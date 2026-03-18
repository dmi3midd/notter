package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNoteNotFound error = errors.New("note not found")

type Note struct {
	Id        string    `json:"id" db:"id"`
	BoardId   *string   `json:"boardId" db:"board_id"`
	UserId    string    `json:"userId" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type NoteDto struct {
	Id        string    `json:"id"`
	BoardId   *string   `json:"boardId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewNoteDto(note *Note) *NoteDto {
	return &NoteDto{
		Id:        note.Id,
		BoardId:   note.BoardId,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

// Need to refactor to unify methods:
// GetNotesByUserId, GetNotesByBoardId, GetStandaloneNotes
type NoteRepository interface {
	GetNote(
		ctx context.Context,
		noteId string,
	) (*Note, error)

	GetNotesByUserId(
		ctx context.Context,
		userId string,
	) ([]Note, error)

	GetNotesByBoardId(
		ctx context.Context,
		boardId string,
	) ([]Note, error)

	GetStandaloneNotes(
		ctx context.Context,
		userId string,
	) ([]Note, error)

	CreateNote(
		ctx context.Context,
		note *Note,
	) error

	UpdateNote(
		ctx context.Context,
		note *Note,
	) error

	DeleteNote(
		ctx context.Context,
		noteId string,
	) error
}
