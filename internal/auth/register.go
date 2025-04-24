package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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

		resp := RegisterResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
