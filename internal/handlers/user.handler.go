package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

		userData, err := h.userService.Registration(
			reqBody.Username,
			reqBody.Email,
			reqBody.Password,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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

		userData, err := h.userService.Login(
			reqBody.Email,
			reqBody.Password,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (h *UserHandler) LogoutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refreshToken")
		if err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}
		refreshToken := cookie.Value
		if err := h.userService.Logout(refreshToken); err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
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
			log.Println(err)
			http.Error(w, "Failed to logout in cookie", http.StatusInternalServerError)
			return
		}
		refreshToken := cookie.Value
		userData, err1 := h.userService.Refresh(refreshToken)
		if err1 != nil {
			log.Println(err1)
			http.Error(w, "Failed to refresh in service", http.StatusInternalServerError)
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
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
