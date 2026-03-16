package services

import (
	"context"
	"fmt"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
)

type NoteService struct {
	repo domain.NoteRepository
}

func NewNoteService(repo domain.NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (ns *NoteService) GetNote(
	ctx context.Context,
	noteId string,
) (*domain.NoteDto, error) {
	op := "note.service-GetNote"
	note, err := ns.repo.GetNote(ctx, noteId)
	noteDto := domain.NewNoteDto(note)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return noteDto, nil
}

func (ns *NoteService) GetNotes(
	ctx context.Context,
	userId string,
) ([]domain.NoteDto, error) {
	op := "note.service-GetNotes"
	notes, err := ns.repo.GetNotes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	notesDto := []domain.NoteDto{}
	for _, note := range notes {
		notesDto = append(notesDto, *domain.NewNoteDto(&note))
	}
	return notesDto, nil
}

func (ns *NoteService) CreateNote(
	ctx context.Context,
	userId string,
	title string,
	content string,
	tags []string,
) error {
	op := "note.service-CreateNote"
	noteId := uuid.NewString()
	if err := ns.repo.CreateNote(ctx, noteId, userId, title, content, tags); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ns *NoteService) UpdateNote(
	ctx context.Context,
	noteId string,
	title string,
	content string,
	tags []string,
) error {
	op := "note.service-UpdateNote"
	if err := ns.repo.UpdateNote(ctx, noteId, title, content, tags); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ns *NoteService) DeleteNote(
	ctx context.Context,
	noteId string,
) error {
	op := "note.service-DeleteNote"
	if err := ns.repo.DeleteNote(ctx, noteId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
