package routers

import (
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewBoardRouter(handler *handlers.BoardHandler) *chi.Mux {
	boardRouter := chi.NewRouter()

	boardRouter.Get("/", handler.GetBoardsHandler())
	boardRouter.Post("/", handler.CreateBoardHandler())
	boardRouter.Put("/{boardId}", handler.UpdateBoard())
	boardRouter.Delete("/{boardId}", handler.DeleteBoard())

	return boardRouter
}
