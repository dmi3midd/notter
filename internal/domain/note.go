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

type NoteRepository interface {
	// If necessary, return ErrNoteNotFound
	GetNote(
		ctx context.Context,
		noteId string,
	) (*Note, error)

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

	// If necessary, return ErrNoteNotFound
	UpdateNote(
		ctx context.Context,
		note *Note,
	) error

	DeleteNote(
		ctx context.Context,
		noteId string,
	) error
}

type NoteService interface {
	GetNote(
		ctx context.Context,
		noteId string,
	) (*NoteDto, error)

	GetNotesByBoardId(
		ctx context.Context,
		boardId string,
	) ([]NoteDto, error)

	GetStandaloneNotes(
		ctx context.Context,
		userId string,
	) ([]NoteDto, error)

	CreateNote(
		ctx context.Context,
		boardId string,
		userId string,
		title string,
		content string,
	) error

	UpdateNote(
		ctx context.Context,
		userId string,
		note NoteDto,
	) error

	DeleteNote(
		ctx context.Context,
		noteId string,
	) error
}
