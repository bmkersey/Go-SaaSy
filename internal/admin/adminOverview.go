package admin

import (
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	httphelpers "github.com/bmkersey/Go-SaaSy/internal/httpHelpers"
)

func AdminOverviewHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgs, _ := store.ListOrganizations(r.Context())
		users, _ := store.ListAllUsers(r.Context())

		orgSafe := make([]map[string]interface{}, 0, len(orgs))
		for _, o := range orgs {
			orgSafe = append(orgSafe, map[string]interface{}{
				"id":            o.ID,
				"name":          o.Name,
				"owner_id":      o.OwnerID,
				"plan":          httphelpers.SafeString(o.Plan),
				"billing_email": httphelpers.SafeString(o.BillingEmail),
				"is_paid":       o.IsPaid,
			})
		}

		httphelpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"orgs":  orgSafe,
			"users": users,
		})
	}
}
