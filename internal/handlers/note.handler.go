package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	noteService domain.NoteService
}

func NewNoteHadnler(noteService domain.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

type CreateOrUpdateNoteRequest struct {
	Title   string
	Content string
}

func (h *NoteHandler) GetNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := chi.URLParam(r, "noteId")
		if noteId == "" {
			http.Error(w, "Note id is required", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		note, err := h.noteService.GetNote(ctx, noteId)
		if err != nil {
			if errors.Is(err, domain.ErrNoteNotFound) {
				http.Error(w, "Note note found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(note); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *NoteHandler) GetStandaloneNotesHandler() http.HandlerFunc {
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

		notes, err := h.noteService.GetStandaloneNotes(ctx, user.Id)
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
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *NoteHandler) GetBoardNotesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boardId := chi.URLParam(r, "boardId")
		if boardId == "" {
			http.Error(w, "Board id is required", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		notes, err := h.noteService.GetNotesByBoardId(ctx, &boardId)
		if err != nil {
			if errors.Is(err, domain.ErrBoardNotFound) {
				http.Error(w, "Board not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *NoteHandler) CreateNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateOrUpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		boardId := chi.URLParam(r, "boardId")
		if boardId == "" {
			http.Error(w, "Board id is required", http.StatusBadRequest)
			return
		}
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

		if err := h.noteService.CreateNote(
			ctx, &boardId,
			user.Id,
			reqBody.Title,
			reqBody.Content,
		); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *NoteHandler) UpdateNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateOrUpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		noteId := chi.URLParam(r, "noteId")
		if noteId == "" {
			http.Error(w, "Note id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.noteService.UpdateNote(ctx, noteId, reqBody.Title, reqBody.Content); err != nil {
			if errors.Is(err, domain.ErrNoteNotFound) {
				http.Error(w, "Note not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *NoteHandler) DeleteNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := chi.URLParam(r, "noteId")
		if noteId == "" {
			http.Error(w, "Note id is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		if err := h.noteService.DeleteNote(ctx, noteId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
