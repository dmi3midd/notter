package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dmi3midd/notter/internal/domain"
)

type NoteHandler struct {
	noteService domain.NoteService
}

func NewNoteHandler(noteService domain.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

type CreateNoteReqBody struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdateNoteReqBody struct {
	NoteId  string   `json:"noteId"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type DeleteNoteReqBody struct {
	NoteId string `json:"noteId"`
}

func (nh *NoteHandler) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userValue := ctx.Value("user")
		if userValue == nil {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		user, ok := userValue.(*domain.User)
		if !ok {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		notes, err := nh.noteService.GetNotes(ctx, user.Id)
		if err != nil {
			log.Printf("ERROR: %v", errors.Unwrap(err))
			http.Error(w, "Invalid request", http.StatusBadRequest)
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

func (nh *NoteHandler) CreateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateNoteReqBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		ctx := r.Context()
		userValue := ctx.Value("user")
		if userValue == nil {
			log.Printf("ERROR in userValue: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		user, ok := userValue.(*domain.User)
		if !ok {
			log.Printf("ERROR in userValue assertion: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		if err := nh.noteService.CreateNote(
			ctx, user.Id,
			reqBody.Title, reqBody.Content, reqBody.Tags,
		); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (nh *NoteHandler) UpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody UpdateNoteReqBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		ctx := r.Context()
		userValue := ctx.Value("user")
		if userValue != nil {
			log.Printf("ERROR in userValue: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}
		_, ok := userValue.(string)
		if !ok {
			log.Printf("ERROR in userValue assertion: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		if err := nh.noteService.UpdateNote(
			ctx,
			reqBody.NoteId, reqBody.Title, reqBody.Content, reqBody.Tags,
		); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (nh *NoteHandler) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody DeleteNoteReqBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		ctx := r.Context()
		userValue := ctx.Value("user")
		if userValue != nil {
			log.Printf("ERROR in userValue: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}
		_, ok := userValue.(string)
		if !ok {
			log.Printf("ERROR in userValue assertion: %v", userValue)
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}

		if err := nh.noteService.DeleteNote(ctx, reqBody.NoteId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
