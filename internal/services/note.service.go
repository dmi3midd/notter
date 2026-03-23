package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
)

type NoteService struct {
	noteStore domain.NoteRepository
	// Need to check if board or user exist
	userStore  domain.UserRepository
	boardStore domain.BoardRepository
}

func NewNoteService(
	noteStore domain.NoteRepository,
	userStore domain.UserRepository,
	boardStore domain.BoardRepository,
) *NoteService {
	return &NoteService{
		noteStore:  noteStore,
		userStore:  userStore,
		boardStore: boardStore,
	}
}

func (s *NoteService) GetNote(
	ctx context.Context,
	noteId string,
) (*domain.NoteDto, error) {
	op := "note.service-GetNote"
	note, err := s.noteStore.GetNote(ctx, noteId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	noteDto := domain.ToNoteDto(note, true)
	return noteDto, nil
}

func (s *NoteService) GetNotesByBoardId(
	ctx context.Context,
	boardId *string,
) ([]domain.NoteDto, error) {
	op := "note.service-GetNotesByBoardId"

	if boardId != nil && *boardId != "" {
		if _, err := s.boardStore.GetBoard(ctx, *boardId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	notes, err := s.noteStore.GetNotesByBoardId(ctx, *boardId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	notesDto := []domain.NoteDto{}
	for _, note := range notes {
		notesDto = append(notesDto, *domain.ToNoteDto(&note, false))
	}
	return notesDto, nil
}

func (s *NoteService) GetStandaloneNotes(
	ctx context.Context,
	userId string,
) ([]domain.NoteDto, error) {
	op := "note.service-GetStandaloneNotes"

	if _, err := s.userStore.GetById(ctx, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	notes, err := s.noteStore.GetStandaloneNotes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	notesDto := []domain.NoteDto{}
	for _, note := range notes {
		notesDto = append(notesDto, *domain.ToNoteDto(&note, false))
	}
	return notesDto, nil
}

// Method to create a note
// Can return domain.ErrUserNotFound or domain.ErrBoardNotFound
func (s *NoteService) CreateNote(
	ctx context.Context,
	boardId *string,
	userId string,
	title string,
	content string,
) error {
	op := "note.service-CreateNote"

	if _, err := s.userStore.GetById(ctx, userId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if boardId != nil && *boardId != "" {
		if _, err := s.boardStore.GetBoard(ctx, *boardId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		boardId = nil
	}

	newNote := domain.Note{
		Id:        uuid.NewString(),
		BoardId:   boardId,
		UserId:    userId,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.noteStore.CreateNote(ctx, &newNote); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *NoteService) UpdateNote(
	ctx context.Context,
	noteId string,
	title string,
	content string,
) error {
	op := "note.service-UpdateNote"
	if err := s.noteStore.UpdateNote(ctx, noteId, title, content, time.Now()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Method to delete the note
func (s *NoteService) DeleteNote(
	ctx context.Context,
	noteId string,
) error {
	op := "note.service-DeleteNote"
	if err := s.noteStore.DeleteNote(ctx, noteId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
