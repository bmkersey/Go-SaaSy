package auth

import (
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/db"
)

func RequireAdmin(store db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !IsAdmin(r, store) {
				http.Error(w, "Admin access only", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
