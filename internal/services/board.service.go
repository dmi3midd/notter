package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/google/uuid"
)

type BoardService struct {
	store     domain.BoardRepository
	userStore domain.UserRepository
}

func NewBoardRepo(
	store domain.BoardRepository,
	userStore domain.UserRepository,
) *BoardService {
	return &BoardService{
		store:     store,
		userStore: userStore,
	}
}

// Method to get boards by user id
// Can return domain.ErrUserNotFound
func (s *BoardService) GetBoards(ctx context.Context, userId string) ([]domain.BoardDto, error) {
	op := "board.service-GetBoards"

	if _, err := s.userStore.GetById(ctx, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	boards, err := s.store.GetBoardsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	boardsDto := []domain.BoardDto{}
	for _, board := range boards {
		boardsDto = append(boardsDto, *domain.NewBoardDto(&board))
	}

	return boardsDto, nil
}

// Method to create board
func (s *BoardService) CreateBoard(
	ctx context.Context,
	userId string,
	title string,
) error {
	op := "board.service-CreateBoards"

	if _, err := s.userStore.GetById(ctx, userId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	board := domain.Board{
		Id:        uuid.NewString(),
		UserId:    userId,
		Title:     title,
		Notes:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.store.CreateBoard(ctx, &board); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Method to update board
// Can return domain.ErrBoardNotFound
func (s *BoardService) UpdateNote(
	ctx context.Context,
	boardId string,
	title string,
) error {
	op := "board.service-UpdateBoard"

	boardCandidate, err := s.store.GetBoard(ctx, boardId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	board := domain.Board{
		Id:        boardId,
		UserId:    boardCandidate.Id,
		Title:     title,
		Notes:     boardCandidate.Notes,
		CreatedAt: boardCandidate.CreatedAt,
		UpdatedAt: time.Now(),
	}
	if err := s.store.UpdateBoard(ctx, &board); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Method to delete board
func (s BoardService) DeleteBoard(ctx context.Context, boardId string) error {
	op := "board.service-DeleteBoard"

	if err := s.store.DeleteBoard(ctx, boardId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
