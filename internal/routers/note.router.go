package routers

import (
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewNoteRouter(handler *handlers.NoteHandler) *chi.Mux {
	noteRouter := chi.NewRouter()

	noteRouter.Get("/board/{boardId}", handler.GetBoardNotesHandler())
	noteRouter.Get("/standalone", handler.GetStandaloneNotesHandler())
	noteRouter.Get("/{noteId}", handler.GetNoteHandler())
	noteRouter.Post("/", handler.CreateNoteHandler())
	noteRouter.Post("/{boardId}", handler.CreateNoteHandler())
	noteRouter.Put("/{noteId}", handler.UpdateNoteHandler())
	noteRouter.Delete("/{noteId}", handler.DeleteNoteHandler())

	return noteRouter
}
