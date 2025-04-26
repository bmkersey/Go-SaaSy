package orgs

import (
	"encoding/json"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetOrgMembers(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		orgIDParam := chi.URLParam(r, "id")
		if orgIDParam == "" {
			http.Error(w, "Organization ID is required", http.StatusBadRequest)
			return
		}

		orgID, err := uuid.Parse(orgIDParam)
		if err != nil {
			http.Error(w, "Invalid organization ID format", http.StatusBadRequest)
			return
		}

		members, err := store.GetOrgMembers(r.Context(), uuid.NullUUID{
			UUID:  orgID,
			Valid: true,
		})
		if err != nil {
			http.Error(w, "Error getting members", http.StatusInternalServerError)
			return
		}

		var resp []db.GetOrgMembersRow
		resp = append(resp, members...)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
