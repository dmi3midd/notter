package routers

import (
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewUserRouter(handler *handlers.UserHandler) *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Post("/registration", handler.RegisterUserHandler())
	userRouter.Post("/login", handler.LoginUserHandler())
	userRouter.Post("/logout", handler.LogoutUserHandler())
	userRouter.Post("/refresh", handler.RefreshTokenHandler())

	return userRouter
}
