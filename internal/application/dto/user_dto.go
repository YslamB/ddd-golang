package dto

import (
	"ddd/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUsersResponse(users []*entity.User, total, page, limit int) UsersResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	return UsersResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}
}
