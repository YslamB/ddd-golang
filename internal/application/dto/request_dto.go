package dto

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type GetUsersRequest struct {
	Page  int `json:"page" validate:"omitempty,min=1"`
	Limit int `json:"limit" validate:"omitempty,min=1,max=100"`
}
