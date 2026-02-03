package http

import (
	"net/http"
	"strings"

	"github.com/maximfill/go-pet-backend/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ParseJWT(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := auth.WithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
