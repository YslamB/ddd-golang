package user

import (
	"context"
	"fmt"

	"gddd/internal/application/user/dtos"
	domain_common "gddd/internal/domain/common"
	domain_user "gddd/internal/domain/user"
	"gddd/internal/infrastructure/logging"
	shared_validation "gddd/internal/shared/validation"

	"github.com/google/uuid"
)

type Service struct {
	userRepo domain_user.Repository
	logger   *logging.Logger
}

func NewService(userRepo domain_user.Repository, log *logging.Logger) *Service {
	return &Service{userRepo: userRepo, logger: log}
}

func (s *Service) RegisterUser(ctx context.Context, req dtos.RegisterUserRequest) (*dtos.UserResponse, error) {
	s.logger.Infof("Attempting to register user with email: %s", req.Email)

	if err := shared_validation.ValidateStruct(req); err != nil {
		s.logger.Warnf("Validation error for user registration: %v", err)
		return nil, domain_common.NewValidationError(err.Error())
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !domain_common.IsNotFoundError(err) {
		s.logger.Errorf("Failed to check existing user by email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		s.logger.Warnf("User with email %s already exists.", req.Email)
		return nil, domain_common.NewConflictError("user with this email already exists")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		s.logger.Errorf("Failed to hash password for email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUserID := domain_user.UserID(uuid.New().String())
	newUser, err := domain_user.NewUser(newUserID, req.Email, hashedPassword)
	if err != nil {
		s.logger.Errorf("Failed to create new user domain entity for email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to create new user domain entity: %w", err)
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Errorf("Failed to save user %s: %v", newUser.ID, err)
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	s.logger.Infof("User %s registered successfully.", newUser.ID)

	return &dtos.UserResponse{
		ID:        string(newUser.ID),
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*dtos.UserResponse, error) {
	s.logger.Infof("Attempting to get user by ID: %s", id)

	userID := domain_user.UserID(id)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("User with ID %s not found.", id)
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
		}
		s.logger.Errorf("Failed to get user by ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	s.logger.Infof("User %s retrieved successfully.", user.ID)
	return &dtos.UserResponse{
		ID:        string(user.ID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, id string, req dtos.UpdateUserRequest) (*dtos.UserResponse, error) {
	s.logger.Infof("Attempting to update user %s", id)

	if err := shared_validation.ValidateStruct(req); err != nil {
		s.logger.Warnf("Validation error for user update %s: %v", id, err)
		return nil, domain_common.NewValidationError(err.Error())
	}

	userID := domain_user.UserID(id)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("User with ID %s not found for update.", id)
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
		}
		s.logger.Errorf("Failed to find user %s for update: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve user for update: %w", err)
	}

	if req.Email != "" && req.Email != user.Email {

		existingUserWithNewEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil && !domain_common.IsNotFoundError(err) {
			s.logger.Errorf("Failed to check for conflicting email %s during update: %v", req.Email, err)
			return nil, fmt.Errorf("failed to check for conflicting email: %w", err)
		}
		if existingUserWithNewEmail != nil && existingUserWithNewEmail.ID != user.ID {
			s.logger.Warnf("New email %s conflicts with existing user %s.", req.Email, existingUserWithNewEmail.ID)
			return nil, domain_common.NewConflictError("email already in use by another user")
		}
		if err := user.ChangeEmail(req.Email); err != nil {
			s.logger.Errorf("Failed to change user %s email to %s: %v", user.ID, req.Email, err)
			return nil, fmt.Errorf("failed to change user email: %w", err)
		}
	}

	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			s.logger.Errorf("Failed to hash new password for user %s: %v", user.ID, err)
			return nil, fmt.Errorf("failed to hash new password: %w", err)
		}
		if err := user.ChangePassword(hashedPassword); err != nil {
			s.logger.Errorf("Failed to change user %s password: %v", user.ID, err)
			return nil, fmt.Errorf("failed to change user password: %w", err)
		}
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Errorf("Failed to save updated user %s: %v", user.ID, err)
		return nil, fmt.Errorf("failed to save updated user: %w", err)
	}

	s.logger.Infof("User %s updated successfully.", user.ID)
	return &dtos.UserResponse{
		ID:        string(user.ID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	s.logger.Infof("Attempting to delete user: %s", id)

	userID := domain_user.UserID(id)
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("User with ID %s not found for deletion.", id)
			return domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
		}
		s.logger.Errorf("Failed to delete user %s: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Infof("User %s deleted successfully.", id)
	return nil
}

func hashPassword(password string) (string, error) {

	return "hashed_" + password, nil
}
