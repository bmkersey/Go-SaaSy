package orgs

import (
	"github.com/bmkersey/Go-SaaSy/internal/db"
)

func HasActivePlan(org db.Organization) bool {
	if !org.IsPaid {
		return false
	}

	if org.Plan.String == "basic" || org.Plan.String == "" {
		return false
	}

	return true
}
