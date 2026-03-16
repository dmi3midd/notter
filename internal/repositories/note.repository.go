package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type NoteRepository struct {
	store *sqlx.DB
}

func NewNoteRepo(store *sqlx.DB) *NoteRepository {
	return &NoteRepository{
		store: store,
	}
}

func (r *NoteRepository) GetNote(
	ctx context.Context,
	noteId string,
) (*domain.Note, error) {
	op := "note.repository-GetNote"
	query := "SELECT * FROM notes WHERE id = $1"
	var note domain.Note
	err := r.store.GetContext(ctx, &note, query, noteId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &note, nil
}

func (r *NoteRepository) GetNotes(
	ctx context.Context,
	userId string,
) ([]domain.Note, error) {
	op := "note.repository-GetNotes"
	query := "SELECT * FROM notes WHERE user_id = $1"
	var notes []domain.Note
	err := r.store.SelectContext(ctx, &notes, query, userId, "")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return notes, nil
}

func (r *NoteRepository) CreateNote(
	ctx context.Context,
	noteId string,
	userId string,
	boardId string,
	title string,
	content string,
	tags []string,
) error {
	op := "note.repository-CreateNote"
	query := `INSERT INTO notes (note_id, user_id, title, content, tags)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := r.store.ExecContext(ctx, query,
		noteId, userId, title, content, pq.Array(tags),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *NoteRepository) UpdateNote(
	ctx context.Context,
	noteId string,
	title string,
	content string,
	tags []string,
) error {
	op := "note.repository-UpdateNote"
	query := `UPDATE notes SET title = $1, content = $2, tags = $3 WHERE note_id = $4`
	var note domain.Note
	err := r.store.GetContext(ctx, &note, query,
		title, content, pq.Array(tags), noteId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNoteNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *NoteRepository) DeleteNote(
	ctx context.Context,
	noteId string,
) error {
	op := "note.repository-DeleteNote"
	query := "DELETE FROM notes WHERE note_id = $1"
	_, err := r.store.ExecContext(ctx, query, noteId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
