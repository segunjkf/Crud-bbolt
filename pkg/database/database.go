package database

import "context"

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Age    int	`json:"age"`
}

// Database represent the  operation that are done on the database.
// This interface abstracts the underlying implemetation.
type Database interface {
	Create(ctx context.Context, data User) error
	GetUser(ctx context.Context, name string) *User
	// Update a a given user (if found) and returns the update user
	Update(ctx context.Context, name string) (User, error)
}