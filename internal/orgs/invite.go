package orgs

import (
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	httphelpers "github.com/bmkersey/Go-SaaSy/internal/httpHelpers"
)

func ValidateInviteHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}

		invite, err := store.GetInviteByToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid or expired invite", http.StatusNotFound)
			return
		}

		org, err := store.GetOrganization(r.Context(), invite.OrgID)
		if err != nil {
			http.Error(w, "Organization not found", http.StatusInternalServerError)
			return
		}

		httphelpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"org_id":   org.ID,
			"org_name": org.Name,
			"email":    invite.Email,
		})
	}
}
