package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	httphelpers "github.com/bmkersey/Go-SaaSy/internal/httpHelpers"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(store db.Store, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			httphelpers.SendError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := store.GetUserByEmail(r.Context(), input.Email)
		if err != nil {
			httphelpers.SendError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err := CheckPassword(user.PasswordHash, input.Password); err != nil {
			httphelpers.SendError(w, http.StatusUnauthorized, "Invalid password/email")
			return
		}

		token, err := GenerateJwt(user.ID.String(), jwtSecret, time.Hour*24)
		if err != nil {
			httphelpers.SendError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
