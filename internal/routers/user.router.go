package routers

import (
	"net/http"

	"github.com/dmi3midd/notter/internal/handlers"
)

func NewUserMux(handler *handlers.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /registration", handler.RegisterUserHandler())
	mux.HandleFunc("POST /login", handler.LoginUserHandler())
	mux.HandleFunc("POST /logout", handler.LogoutUserHandler())
	mux.HandleFunc("POST /refresh", handler.RefreshTokenHandler())

	return mux
}
