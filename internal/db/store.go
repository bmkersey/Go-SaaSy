package db

import (
	"context"

	"github.com/google/uuid"
)

// Store is an interface that wraps all DB methods we use in handlers.
type Store interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	UpdateUserOrg(ctx context.Context, arg UpdateUserOrgParams) error
	CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (CreateOrganizationRow, error)
	GetOrganization(ctx context.Context, id uuid.UUID) (Organization, error)
	GetOrgMembers(ctx context.Context, organizationID uuid.NullUUID) ([]GetOrgMembersRow, error)
	UpdateOrganizationBilling(ctx context.Context, arg UpdateOrganizationBillingParams) error
	CreatePasswordReset(ctx context.Context, arg CreatePasswordResetParams) error
	GetPasswordResetByToken(ctx context.Context, token string) (PasswordReset, error)
	DeletePasswordReset(ctx context.Context, token string) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	ListAllUsers(ctx context.Context) ([]ListAllUsersRow, error)
	ListOrganizations(ctx context.Context) ([]ListOrganizationsRow, error)
	SetUserAdmin(ctx context.Context, email string) error
}

// SQLStore implements the Store interface using sqlc-generated Queries.
type SQLStore struct {
	*Queries
}

// NewStore returns a new SQLStore.
func NewStore(q *Queries) Store {
	return &SQLStore{q}
}
