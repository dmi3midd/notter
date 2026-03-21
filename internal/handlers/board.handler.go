package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/go-chi/chi/v5"
)

type BoardHandler struct {
	service domain.BoardService
}

func NewBoardHandler(boardService domain.BoardService) *BoardHandler {
	return &BoardHandler{
		service: boardService,
	}
}

type CreateOrUpdateBoardRequest struct {
	Title string
}

func (h *BoardHandler) GetBoardsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userFromCtx := ctx.Value("user")
		if userFromCtx == nil {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		user, ok := userFromCtx.(domain.User)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		boards, err := h.service.GetBoards(ctx, user.Id)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				http.Error(w, "User doesn't exist", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(boards); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *BoardHandler) CreateBoardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateOrUpdateBoardRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		userFromCtx := ctx.Value("user")
		if userFromCtx == nil {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		user, ok := userFromCtx.(domain.User)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := h.service.CreateBoard(ctx, user.Id, reqBody.Title); err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				http.Error(w, "User doesn't exist", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *BoardHandler) UpdateBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateOrUpdateBoardRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		boardId := chi.URLParam(r, "boardId")
		if boardId == "" {
			http.Error(w, "Board id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.service.UpdateBoard(ctx, boardId, reqBody.Title); err != nil {
			if errors.Is(err, domain.ErrBoardNotFound) {
				http.Error(w, "Board not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *BoardHandler) DeleteBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boardId := chi.URLParam(r, "boardId")
		if boardId == "" {
			http.Error(w, "Board id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.service.DeleteBoard(ctx, boardId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
