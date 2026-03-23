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
	BoardId   *string   `json:"boardId,omitempty"`
	Title     string    `json:"title"`
	Content   *string   `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToNoteDto(note *Note, includeContent bool) *NoteDto {
	dto := &NoteDto{
		Id:        note.Id,
		BoardId:   note.BoardId,
		Title:     note.Title,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	if includeContent {
		dto.Content = &note.Content
	}

	return dto
}

type NoteRepository interface {
	// GetNote retrieves a Note entity by its id.
	// It returns ErrNoteNotFound if no note are found.
	GetNote(ctx context.Context, noteId string) (*Note, error)
	// GetNotesByBoardId retrieves a Note slice by board id.
	// It returns an empty slice if no notes are found.
	GetNotesByBoardId(ctx context.Context, boardId string) ([]Note, error)
	// GetStandaloneNotes retrieves a standalone Note slice by userId.
	// It returns an empty slice if no notes are found.
	GetStandaloneNotes(ctx context.Context, userId string) ([]Note, error)
	// CreateNote inserts a new note and increments the note count on the associated board.
	// It returns domain.ErrBoardNotFound if the specified board does not exist.
	CreateNote(ctx context.Context, note *Note) error
	// UpdateNote modifies the title and content of an existing note.
	// It returns domain.ErrNoteNotFound if no note is found.
	UpdateNote(ctx context.Context, noteId string, title string, content string, updateAt time.Time) error
	// DeleteNote removes the note.
	DeleteNote(ctx context.Context, noteId string) error
}

type NoteService interface {
	// GetNote finds note and returns NoteDto struct.
	// It returns ErrNoteNotFound if no note are found.
	GetNote(ctx context.Context, noteId string) (*NoteDto, error)
	// GetNotesByBoardId finds a notes by board id and returns NoteDto slice.
	// It returns an empty slice if no notes are found.
	// It returns ErrBoardNotFound if no board are found.
	GetNotesByBoardId(ctx context.Context, boardId *string) ([]NoteDto, error)
	// GetStandaloneNotes finds a notes by user id and returns NoteDto slice.
	// It returns ErrUserNotFound if no user are found.
	// It returns an empty slice if no notes are found.
	GetStandaloneNotes(ctx context.Context, userId string) ([]NoteDto, error)
	// CreateNote creates a new note.
	// It returns ErrUserNotFound if no user are found.
	// It returns domain.ErrBoardNotFound if the specified board does not exist.
	CreateNote(ctx context.Context, boardId *string, userId string, title string, content string) error
	// UpdateNote modifies the existing note.
	// It returns domain.ErrNoteNotFound if no note is found.
	UpdateNote(ctx context.Context, noteId string, title string, content string) error
	// DeleteNote removes the note.
	DeleteNote(ctx context.Context, noteId string) error
}
