package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/models"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"gorm.io/gorm"
)

// UserRepository interface defines user-specific repository methods
type UserRepository interface {
	Repository[models.User]
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByRefreshToken(ctx context.Context, token string) (*models.User, error)
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error
	UpdateStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	SearchUsers(ctx context.Context, query string, page, pageSize int) ([]models.User, int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	*BaseRepository[models.User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[models.User](db),
	}
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByRefreshToken finds a user by refresh token
func (r *userRepository) FindByRefreshToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("refresh_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateRefreshToken updates the user's refresh token
func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", token).Error
}

// UpdateLastLogin updates the user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()")).Error
}

// VerifyEmail marks the user's email as verified
func (r *userRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("email_verified_at", gorm.Expr("NOW()")).Error
}

// UpdatePassword updates the user's password
func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Note: The password should be hashed before calling this method, or use the BeforeUpdate hook
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", password).Error
}

// UpdateStatus updates the user's status
func (r *userRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status models.UserStatus) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}

// UpdateRole updates the user's role
func (r *userRepository) UpdateRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// ExistsByEmail checks if a user exists with the given email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// SearchUsers searches for users by name or email
func (r *userRepository) SearchUsers(ctx context.Context, query string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	searchQuery := "%" + query + "%"

	// Count total
	err := r.DB.WithContext(ctx).Model(&models.User{}).
		Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", searchQuery, searchQuery, searchQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = r.DB.WithContext(ctx).
		Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// FindByID overrides base to include error handling
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
