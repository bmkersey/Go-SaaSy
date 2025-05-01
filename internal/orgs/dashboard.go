package orgs

import (
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	httphelpers "github.com/bmkersey/Go-SaaSy/internal/httpHelpers"
	"github.com/google/uuid"
)

func OrgDashboardHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgIDStr, ok := r.Context().Value(OrgIDKey).(string)
		if !ok || orgIDStr == "" {
			http.Error(w, "Missing org ID", http.StatusUnauthorized)
			return
		}

		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			http.Error(w, "Invalid org ID", http.StatusBadRequest)
			return
		}

		org, err := store.GetOrganization(r.Context(), orgID)
		if err != nil {
			http.Error(w, "Error loading org", http.StatusInternalServerError)
			return
		}

		members, err := store.GetOrgMembers(r.Context(), uuidNull(orgID))
		if err != nil {
			http.Error(w, "Error fetching members", http.StatusInternalServerError)
			return
		}

		httphelpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"org": map[string]interface{}{
				"name":          org.Name,
				"plan":          org.Plan.String,
				"billing_email": org.BillingEmail.String,
				"is_paid":       org.IsPaid,
			},
			"members": members,
		})
	}
}

func uuidNull(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: id, Valid: true}
}
