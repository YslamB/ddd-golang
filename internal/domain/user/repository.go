package user

import "context"

// Repository defines the interface for User persistence operations.
// This interface is part of the Domain Layer, meaning the Domain
// specifies *what* persistence operations are needed, not *how* they are performed.
type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id UserID) error
}
