package user

import (
	"errors"
	"time"
)

// UserID is a Value Object for user identification.
// It's immutable and its value determines its equality.
type UserID string

// User is an Entity/Aggregate Root. It has a unique identity (ID)
// and encapsulates business logic related to user state.
type User struct {
	ID        UserID
	Email     string
	Password  string // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity.
// It ensures that essential invariants (like non-empty email/password) are met at creation.
func NewUser(id UserID, email, hashedPassword string) (*User, error) {
	if email == "" || hashedPassword == "" {
		return nil, errors.New("email and password cannot be empty")
	}
	return &User{
		ID:        id,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ChangeEmail updates the user's email address. This method encapsulates the business rule
// that an email cannot be empty and ensures the 'UpdatedAt' timestamp is updated.
func (u *User) ChangeEmail(newEmail string) error {
	if newEmail == "" {
		return errors.New("new email cannot be empty")
	}
	u.Email = newEmail
	u.UpdatedAt = time.Now()
	return nil
}

// ChangePassword updates the user's hashed password. Similar to ChangeEmail, it enforces
// the non-empty rule and updates the timestamp. The actual hashing logic is outside the domain,
// as the domain only cares about the *hashed* password.
func (u *User) ChangePassword(newHashedPassword string) error {
	if newHashedPassword == "" {
		return errors.New("new password cannot be empty")
	}
	u.Password = newHashedPassword
	u.UpdatedAt = time.Now()
	return nil
}
