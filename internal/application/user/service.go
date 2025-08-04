package user

import (
	"context"
	"fmt"

	"gddd/internal/application/user/dtos" // Import DTOs
	domain_common "gddd/internal/domain/common"
	domain_user "gddd/internal/domain/user" // Import domain package
	infra_logging "gddd/internal/infrastructure/logging"
	shared_validation "gddd/internal/shared/validation"

	"github.com/google/uuid"
)

// UserRepository defines the interface for the user repository,
// used by the application service. This is a copy of the domain interface
// to explicitly define application layer dependencies.
// type UserRepository interface {
// 	Save(ctx context.Context, user *domain_user.User) error
// 	FindByID(ctx context.Context, id domain_user.UserID) (*domain_user.User, error)
// 	FindByEmail(ctx context.Context, email string) (*domain_user.User, error)
// 	Delete(ctx context.Context, id domain_user.UserID) error
// }

// Service defines the application service for user-related operations.
type Service struct {
	userRepo domain_user.Repository
	logger   *infra_logging.Logger
}

// NewService creates a new user application service.
func NewService(userRepo domain_user.Repository, logger *infra_logging.Logger) *Service {
	return &Service{userRepo: userRepo, logger: logger}
}

// RegisterUser handles the user registration use case.
func (s *Service) RegisterUser(ctx context.Context, req dtos.RegisterUserRequest) (*dtos.UserResponse, error) {
	s.logger.Infof("Attempting to register user with email: %s", req.Email)

	// 1. Validate input using shared validation utility
	if err := shared_validation.ValidateStruct(req); err != nil {
		s.logger.Warnf("Validation error for user registration: %v", err)
		return nil, domain_common.NewValidationError(err.Error())
	}

	// 2. Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !domain_common.IsNotFoundError(err) {
		s.logger.Errorf("Failed to check existing user by email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		s.logger.Warnf("User with email %s already exists.", req.Email)
		return nil, domain_common.NewConflictError("user with this email already exists")
	}

	// 3. Hash password (infrastructure concern, but done here for simplicity)
	hashedPassword, err := hashPassword(req.Password) // Placeholder for actual hashing
	if err != nil {
		s.logger.Errorf("Failed to hash password for email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Create domain entity
	newUserID := domain_user.UserID(uuid.New().String())
	newUser, err := domain_user.NewUser(newUserID, req.Email, hashedPassword)
	if err != nil {
		s.logger.Errorf("Failed to create new user domain entity for email %s: %v", req.Email, err)
		return nil, fmt.Errorf("failed to create new user domain entity: %w", err)
	}

	// 5. Persist domain entity
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Errorf("Failed to save user %s: %v", newUser.ID, err)
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	s.logger.Infof("User %s registered successfully.", newUser.ID)

	// 6. Return response DTO
	return &dtos.UserResponse{
		ID:        string(newUser.ID),
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// GetUserByID handles retrieving a user by ID use case.
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

// UpdateUser handles updating a user's details use case.
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
		// Check if new email conflicts with existing user
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

// DeleteUser handles deleting a user use case.
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

// Placeholder for a password hashing function.
// In a real application, use bcrypt or similar for secure hashing.
func hashPassword(password string) (string, error) {
	// For production, use a strong hashing library like golang.org/x/crypto/bcrypt
	return "hashed_" + password, nil
}
