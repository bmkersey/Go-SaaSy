package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/google/uuid"
)

func ForgotPasswordHandler(store db.Store, email *email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddr := r.FormValue("email")

		user, err := store.GetUserByEmail(r.Context(), emailAddr)
		if err != nil {
			http.Error(w, "Error finding user", http.StatusNotFound)
			return
		}

		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		token := base64.URLEncoding.EncodeToString(tokenBytes)

		expiry := time.Now().Add(10 * time.Minute)
		err = store.CreatePasswordReset(r.Context(), db.CreatePasswordResetParams{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: expiry,
		})
		if err != nil {
			http.Error(w, "Error saving new pw reset", http.StatusInternalServerError)
			return
		}

		resetLink := fmt.Sprintf("https://localhost:8080/reset-password?token=%s", token)

		go func() {
			data := struct {
				Name string
				Link string
			}{
				Name: user.Email,
				Link: resetLink,
			}
			err := email.SendTemplate(emailAddr, "Reset Your SaaSy Password", "reset.html.tmpl", data)
			if err != nil {
				fmt.Printf("Error sending reset email: %v\n", err)
			}
		}()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset link sent."))

	}
}
