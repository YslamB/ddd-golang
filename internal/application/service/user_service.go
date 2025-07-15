package service

import (
	"context"
	"ddd/internal/application/dto"
	"ddd/internal/domain/entity"
	"ddd/internal/domain/repository"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {

	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user, err := entity.NewUser(req.Name, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *UserService) GetUsers(ctx context.Context, req dto.GetUsersRequest) (*dto.UsersResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	users, err := s.repo.GetAll(ctx, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	response := dto.ToUsersResponse(users, total, req.Page, req.Limit)
	return &response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.repo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser != nil {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	err = user.Update(req.Name, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
