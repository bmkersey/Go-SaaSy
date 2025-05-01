package auth

import (
	"encoding/json"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
)

type MeResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func MeHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDUUID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Unable to parse userID", http.StatusInternalServerError)
			return
		}

		user, err := store.GetUserByID(r.Context(), userIDUUID)
		if err != nil {
			http.Error(w, "Unable to get user from ID", http.StatusInternalServerError)
			return
		}

		resp := MeResponse{
			ID:      user.ID.String(),
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
