package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/repository"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
)

// UpdateUserRequest represents the update user request
type UpdateUserRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

// UserService interface defines user service methods
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, page, pageSize int) ([]models.User, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error
	UpdateRole(ctx context.Context, id uuid.UUID, role models.UserRole) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll retrieves all users with pagination
func (s *userService) GetAll(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		logger.Error("Failed to get users", logger.Err(err))
		return nil, 0, apperrors.ErrInternal
	}
	return users, total, nil
}

// Update updates a user
func (s *userService) Update(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user", logger.Err(err))
		return nil, apperrors.ErrInternal
	}

	return user, nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", logger.Err(err))
		return apperrors.ErrInternal
	}

	return nil
}

// Search searches for users
func (s *userService) Search(ctx context.Context, query string, page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.SearchUsers(ctx, query, page, pageSize)
	if err != nil {
		logger.Error("Failed to search users", logger.Err(err))
		return nil, 0, apperrors.ErrInternal
	}

	return users, total, nil
}

// UpdateStatus updates a user's status
func (s *userService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdateStatus(ctx, id, status); err != nil {
		logger.Error("Failed to update user status", logger.Err(err))
		return apperrors.ErrInternal
	}

	return nil
}

// UpdateRole updates a user's role
func (s *userService) UpdateRole(ctx context.Context, id uuid.UUID, role models.UserRole) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdateRole(ctx, id, role); err != nil {
		logger.Error("Failed to update user role", logger.Err(err))
		return apperrors.ErrInternal
	}

	return nil
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}
