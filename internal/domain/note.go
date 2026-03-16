package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNoteNotFound error = errors.New("note not found")

type Note struct {
	NoteId    string    `json:"noteId" db:"note_id"`
	UserId    string    `json:"userId" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type NoteDto struct {
	NoteId    string    `json:"noteId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNoteDto(note *Note) *NoteDto {
	return &NoteDto{
		NoteId:    note.NoteId,
		Title:     note.Title,
		Content:   note.Content,
		Tags:      note.Tags,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

type NoteRepository interface {
	GetNote(
		ctx context.Context,
		noteId string,
	) (*Note, error)

	GetNotes(
		ctx context.Context,
		userId string,
	) ([]Note, error)

	CreateNote(
		ctx context.Context,
		noteId string,
		userId string,
		title string,
		content string,
		tags []string,
	) error

	UpdateNote(
		ctx context.Context,
		noteId string,
		title string,
		content string,
		tags []string,
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

	GetNotes(
		ctx context.Context,
		userId string,
	) ([]NoteDto, error)

	CreateNote(
		ctx context.Context,
		userId string,
		title string,
		content string,
		tags []string,
	) error

	UpdateNote(
		ctx context.Context,
		noteId string,
		title string,
		content string,
		tags []string,
	) error

	DeleteNote(
		ctx context.Context,
		noteId string,
	) error
}
