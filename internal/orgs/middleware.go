package orgs

import (
	"context"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
)

type contextKey string

const orgIDKey contextKey = "orgID"

func OrgMiddleware(store db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := auth.GetUserID(r)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			userUUID, err := uuid.Parse(userID)
			if err != nil {
				http.Error(w, "Error parsing userID to UUID", http.StatusInternalServerError)
				return
			}
			user, err := store.GetUserByID(r.Context(), userUUID)
			if err != nil || !user.OrganizationID.Valid {
				http.Error(w, "User does not belong to an organization", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), orgIDKey, user.OrganizationID.UUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetOrgID(r *http.Request) (string, bool) {
	val := r.Context().Value(orgIDKey)
	id, ok := val.(uuid.UUID)
	if !ok {
		return "", false
	}
	return id.String(), ok
}

func RequirePaidPlan(store db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgID, ok := GetOrgID(r)
			if !ok {
				http.Error(w, "Organization not found", http.StatusUnauthorized)
				return
			}

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusBadRequest)
				return
			}

			org, err := store.GetOrganization(r.Context(), orgUUID)
			if err != nil {
				http.Error(w, "Could not find organization", http.StatusInternalServerError)
				return
			}

			if !HasActivePlan(org) {
				http.Error(w, "Subscription required", http.StatusPaymentRequired)
			}

			next.ServeHTTP(w, r)
		})
	}
}
