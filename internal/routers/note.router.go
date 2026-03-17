package routers

import (
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewNoteRouter(handler *handlers.NoteHandler) *chi.Mux {
	noteRouter := chi.NewRouter()

	noteRouter.Get("/notes", handler.GetNotes())
	noteRouter.Post("/note", handler.CreateNote())
	noteRouter.Put("/note", handler.UpdateNote())
	noteRouter.Delete("/note", handler.DeleteNote())

	return noteRouter
}
