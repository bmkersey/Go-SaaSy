package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	httphelpers "github.com/bmkersey/Go-SaaSy/internal/httpHelpers"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, fmt.Sprintf("error registering user: %s\n", err), http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" {
			http.Error(w, "Must supply Email AND Password", http.StatusBadRequest)
			return
		}

		hashed, err := HashPassword(input.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error hashing password: %s\n", err), http.StatusBadRequest)
			return
		}

		user, err := store.CreateUser(r.Context(), db.CreateUserParams{
			ID:           uuid.New(),
			Email:        input.Email,
			PasswordHash: hashed,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %s\n", err), http.StatusBadRequest)
			return
		}

		token, err := GenerateJwt(user.ID.String(), "saasy", time.Hour*24)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		sender := email.NewEmailSender()

		go func() {
			data := struct {
				Name string
			}{
				Name: "Test User",
			}

			if err := sender.SendTemplate(user.Email, "Welcome to SaaSy!", "welcome.html.tmpl", data); err != nil {
				log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
			}
		}()

		httphelpers.RespondWithJSON(w, http.StatusCreated, map[string]string{
			"token": token,
		})
	}
}
