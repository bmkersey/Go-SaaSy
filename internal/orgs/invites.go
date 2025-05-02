package orgs

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/google/uuid"
)

type InviteRequest struct {
	Email string `json:"email"`
}

func CreateInviteHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input InviteRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Email == "" {
			http.Error(w, "Missing or invalid email", http.StatusBadRequest)
			return
		}

		orgIDstr := r.Context().Value(OrgIDKey).(string)
		orgID, _ := uuid.Parse(orgIDstr)

		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		token := base64.URLEncoding.EncodeToString(tokenBytes)

		invite := db.CreateInviteParams{
			ID:        uuid.New(),
			OrgID:     orgID,
			Token:     token,
			Email:     input.Email,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		if err := store.CreateInvite(r.Context(), invite); err != nil {
			http.Error(w, "Failed to create invite", http.StatusInternalServerError)
			return
		}

		link := "http://localhost:3000/invite?token=" + token

		sender := email.NewEmailSender()
		go func() {
			data := struct {
				InviteLink string
			}{InviteLink: link}

			err := sender.SendTemplate(input.Email, "You're invited to SaaSy", "invite.html.tmpl", data)
			if err != nil {
				log.Printf("❌ Failed to send invite email: %v", err)
			}
		}()

		w.WriteHeader(http.StatusCreated)
	}
}
