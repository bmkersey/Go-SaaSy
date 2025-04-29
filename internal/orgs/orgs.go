package orgs

import (
	"encoding/json"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
)

type CreateOrgInput struct {
	Name string `json:"name"`
}

type CreateOrgResponse struct {
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	OwnerID        uuid.UUID `json:"owner_id"`
}

func CreateOrganizationHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDUUID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Unable to parse userID", http.StatusInternalServerError)
			return
		}

		var input CreateOrgInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		newOrgID := uuid.New()
		org, err := store.CreateOrganization(r.Context(), db.CreateOrganizationParams{
			ID:      newOrgID,
			Name:    input.Name,
			OwnerID: userIDUUID,
		})
		if err != nil {
			http.Error(w, "Error creating organization", http.StatusInternalServerError)
			return
		}

		err = store.UpdateUserOrg(r.Context(), db.UpdateUserOrgParams{
			ID: userIDUUID,
			OrganizationID: uuid.NullUUID{
				UUID:  org.ID,
				Valid: true,
			},
		})
		if err != nil {
			http.Error(w, "Error updating User org", http.StatusInternalServerError)
			return
		}

		resp := CreateOrgResponse{
			OrganizationID: org.ID.String(),
			Name:           org.Name,
			OwnerID:        org.OwnerID,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
