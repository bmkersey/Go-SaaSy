package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
	Email       string `json:"email"`
}

func ResetPasswordHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ResetPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}

		reset, err := store.GetPasswordResetByToken(r.Context(), req.Token)
		if err != nil {
			http.Error(w, "Unable to find pw reset or invalid token", http.StatusBadRequest)
			return
		}

		if time.Now().After(reset.ExpiresAt) {
			_ = store.DeletePasswordReset(r.Context(), req.Token)
			http.Error(w, "Token has expired", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		err = store.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
			PasswordHash: string(hashedPassword),
			Email:        req.Email,
		})
		if err != nil {
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		_ = store.DeletePasswordReset(r.Context(), req.Token)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password successfully reset"))
	}
}
