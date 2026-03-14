package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dmi3midd/notter/internal/domain"
)

func Authorization(
	tokenService domain.TokenService,
	userRepository domain.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			authHeader := r.Header.Get("Authorization")
			if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
				token = after
			} else {
				token = r.URL.Query().Get("token")
			}

			if token == "" {
				http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
				return
			}

			payload := tokenService.ValidateAccessToken(token)
			if payload == nil {
				http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
				return
			}

			user, err := userRepository.GetById(r.Context(), payload.Id)
			if err != nil {
				if errors.Is(err, domain.ErrUserNotFound) {
					http.Error(w, domain.ErrUnuthorized.Error(), http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
