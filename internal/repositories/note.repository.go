package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/jmoiron/sqlx"
)

type NoteRepository struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) *NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

// Can return domain.ErrNoteNotFound
func (r *NoteRepository) GetNote(
	ctx context.Context,
	noteId string,
) (*domain.Note, error) {
	op := "note.repository-GetNote"
	query := "SELECT * FROM notes WHERE id = $1"
	var note domain.Note
	if err := r.db.GetContext(ctx, &note, query, noteId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &note, nil
}

func (r *NoteRepository) GetNotesByBoardId(
	ctx context.Context,
	boardId string,
) ([]domain.Note, error) {
	op := "note.repository-GetNotesByBoardId"
	query := "SELECT * FROM notes WHERE board_id = $1"
	var notes []domain.Note
	if err := r.db.SelectContext(ctx, &notes, query, boardId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return notes, nil
}

func (r *NoteRepository) GetStandaloneNotes(
	ctx context.Context,
	userId string,
) ([]domain.Note, error) {
	op := "note.repository-GetStandaloneNotes"
	query := "SELECT * FROM notes WHERE user_id = $1 AND board_id IS NULL"
	var notes []domain.Note
	if err := r.db.SelectContext(ctx, &notes, query, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

func (r *NoteRepository) CreateNote(
	ctx context.Context,
	note *domain.Note,
) error {
	op := "note.repository-CreateNote"
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback()

	queryToCreateNote := `
	INSERT INTO notes (id, board_id, user_id, title, content, created_at, update_at)
    VALUES (:id, :board_id, :user_id, :title, :content, :created_at, :update_at)
	`
	if _, err := r.db.NamedExecContext(ctx, queryToCreateNote, note); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	queryToUpdateBoard := `UPDATE boards SET notes = notes + 1 WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, queryToUpdateBoard, note.BoardId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Can return domain.ErrNoteNotFound
func (r *NoteRepository) UpdateNote(
	ctx context.Context,
	note *domain.Note,
) error {
	op := "note.repository-UpdateNote"
	query := `UPDATE notes SET
							board_id = :board_id,
							title = :title,
							content = :content,
							updated_at = :updated_at
						WHERE id = :id
	`
	if _, err := r.db.NamedExecContext(ctx, query, note); err != nil {
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
	if _, err := r.db.ExecContext(ctx, query, noteId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
