package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
)

type NoteService struct {
	store domain.NoteRepository
	// Need to check if board or user exist
	userStore  domain.UserRepository
	boardStore domain.BoardRepository
}

func NewNoteService(
	store domain.NoteRepository,
	userStore domain.UserRepository,
	boardStore domain.BoardRepository,
) *NoteService {
	return &NoteService{
		store:      store,
		userStore:  userStore,
		boardStore: boardStore,
	}
}

// Method to get the note by it id
// Can return domain.ErrNoteNotFound
func (s *NoteService) GetNote(
	ctx context.Context,
	noteId string,
) (*domain.NoteDto, error) {
	op := "note.service-GetNote"
	note, err := s.store.GetNote(ctx, noteId)
	noteDto := domain.ToNoteDto(note, true)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return noteDto, nil
}

// Method to get notes by board id
// Can return domain.ErrBoardNotFound
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

	notes, err := s.store.GetNotesByBoardId(ctx, *boardId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	notesDto := []domain.NoteDto{}
	for _, note := range notes {
		notesDto = append(notesDto, *domain.ToNoteDto(&note, false))
	}
	return notesDto, nil
}

// Method to get standalone notes by user id
// Can return domain.ErrUserNotFound
func (s *NoteService) GetStandaloneNotes(
	ctx context.Context,
	userId string,
) ([]domain.NoteDto, error) {
	op := "note.service-GetStandaloneNotes"

	if _, err := s.userStore.GetById(ctx, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	notes, err := s.store.GetStandaloneNotes(ctx, userId)
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
	if err := s.store.CreateNote(ctx, &newNote); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Method to update the note
// Can return domain.ErrUserNotFound or domain.ErrBoardNotFound
// func (s *NoteService) UpdateNote(
// 	ctx context.Context,
// 	userId string,
// 	note domain.NoteDto,
// ) error {
// 	op := "note.service-UpdateNote"

// 	if _, err := s.userStore.GetById(ctx, userId); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	if note.BoardId != nil && *note.BoardId != "" {
// 		if _, err := s.boardStore.GetBoard(ctx, *note.BoardId); err != nil {
// 			return fmt.Errorf("%s: %w", op, err)
// 		}
// 	}

// 	newNote := domain.Note{
// 		Id:        note.Id,
// 		BoardId:   note.BoardId,
// 		UserId:    userId,
// 		Title:     note.Title,
// 		Content:   *note.Content,
// 		CreatedAt: note.CreatedAt,
// 		UpdatedAt: time.Now(),
// 	}
// 	if err := s.store.UpdateNote(ctx, &newNote); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	return nil
// }

// Method to delete the note
func (s *NoteService) DeleteNote(
	ctx context.Context,
	noteId string,
) error {
	op := "note.service-DeleteNote"
	if err := s.store.DeleteNote(ctx, noteId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
