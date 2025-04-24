package db

import "context"

// Store is an interface that wraps all DB methods we use in handlers.
type Store interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
}

// SQLStore implements the Store interface using sqlc-generated Queries.
type SQLStore struct {
	*Queries
}

// NewStore returns a new SQLStore.
func NewStore(q *Queries) Store {
	return &SQLStore{q}
}
