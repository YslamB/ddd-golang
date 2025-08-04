package user

import (
	"errors"
	"time"
)

type UserID string

type User struct {
	ID        UserID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

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

func (u *User) ChangeEmail(newEmail string) error {
	if newEmail == "" {
		return errors.New("new email cannot be empty")
	}
	u.Email = newEmail
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ChangePassword(newHashedPassword string) error {
	if newHashedPassword == "" {
		return errors.New("new password cannot be empty")
	}
	u.Password = newHashedPassword
	u.UpdatedAt = time.Now()
	return nil
}
