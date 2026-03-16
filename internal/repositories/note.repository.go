package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type noteRow struct {
	NoteId    string         `db:"note_id"`
	UserId    string         `db:"user_id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Tags      pq.StringArray `db:"tags"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (r *noteRow) toDomain() domain.Note {
	return domain.Note{
		NoteId:    r.NoteId,
		UserId:    r.UserId,
		Title:     r.Title,
		Content:   r.Content,
		Tags:      []string(r.Tags),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

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
	var row noteRow
	err := r.store.GetContext(ctx, &row, query, noteId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	note := row.toDomain()
	return &note, nil
}

func (r *NoteRepository) GetNotes(
	ctx context.Context,
	userId string,
) ([]domain.Note, error) {
	op := "note.repository-GetNotes"
	query := "SELECT * FROM notes WHERE user_id = $1"
	var rows []noteRow
	err := r.store.SelectContext(ctx, &rows, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	notes := make([]domain.Note, len(rows))
	for i, row := range rows {
		notes[i] = row.toDomain()
	}
	
	return notes, nil
}

func (r *NoteRepository) CreateNote(
	ctx context.Context,
	noteId string,
	userId string,
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
