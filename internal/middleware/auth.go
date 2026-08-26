package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/token"
)

type contextKey string

const USER_ID_KEY contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public routes skip auth.
		if isPublic(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authorization := r.Header.Get("Authorization")

		splited := strings.SplitN(authorization, " ", 2)
		if len(splited) != 2 || splited[0] != "Bearer" {
			network.Error(w, http.StatusUnauthorized)
			return
		}

		// This returns user id.
		userId, err := token.Validate(splited[1])
		if err != nil {
			network.Error(w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), USER_ID_KEY, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func isPublic(path string) bool {
	public := []string{"/api/auth/login"}
	return slices.Contains(public, path)
}
