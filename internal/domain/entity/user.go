package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email string) (*User, error) {
	if err := validateUserData(name, email); err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) Update(name, email string) error {
	if name == "" && email == "" {
		return errors.New("at least one field must be provided for update")
	}

	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()

	return validateUserData(u.Name, u.Email)
}

func validateUserData(name, email string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(email) < 5 || !contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
