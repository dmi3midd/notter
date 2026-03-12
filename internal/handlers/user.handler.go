package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dmi3midd/notter/internal/domain"
	"github.com/dmi3midd/notter/internal/services"
)

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody RegistrationRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		userData, err := h.userService.Registration(
			ctx,
			reqBody.Username,
			reqBody.Email,
			reqBody.Password,
		)

		if err != nil {
			log.Printf("ERROR: %v", errors.Unwrap(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    userData.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			HttpOnly: true,
			Path:     "/",
			// Secure: true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *UserHandler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		userData, err := h.userService.Login(
			ctx,
			reqBody.Email,
			reqBody.Password,
		)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				log.Printf("ERROR: %v", errors.Unwrap(err))
				http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    userData.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			HttpOnly: true,
			Path:     "/",
			// Secure: true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *UserHandler) LogoutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refreshToken")
		if err != nil {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}
		refreshToken := cookie.Value
		ctx := r.Context()
		if err := h.userService.Logout(ctx, refreshToken); err != nil {
			log.Printf("ERROR: %v", errors.Unwrap(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) RefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refreshToken")
		if err != nil {
			http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
			return
		}
		refreshToken := cookie.Value
		ctx := r.Context()
		userData, err := h.userService.Refresh(ctx, refreshToken)
		if err != nil {
			log.Printf("ERROR: %v", errors.Unwrap(err))
			if errors.Is(err, domain.ErrUnuthorized) {
				http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    userData.RefreshToken,
			MaxAge:   30 * 24 * 60 * 60,
			HttpOnly: true,
			Path:     "/",
			// Secure: true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
