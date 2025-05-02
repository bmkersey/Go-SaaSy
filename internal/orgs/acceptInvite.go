package orgs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
)

func AcceptInviteHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}

		userID, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		invite, err := store.GetInviteByToken(r.Context(), req.Token)
		if err != nil {
			http.Error(w, "Invalid or expired invite", http.StatusBadRequest)
			return
		}

		err = store.UpdateUserOrg(r.Context(), db.UpdateUserOrgParams{
			ID:             userUUID,
			OrganizationID: uuid.NullUUID{UUID: invite.OrgID, Valid: true},
		})
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		err = store.MarkInviteUsed(r.Context(), invite.ID)
		if err != nil {
			log.Printf("Warning: failed to mark invite used: %v", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
