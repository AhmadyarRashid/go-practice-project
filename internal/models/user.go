package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole represents user roles
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleModerator UserRole = "moderator"
)

// UserStatus represents user account status
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
	StatusPending  UserStatus = "pending"
)

// User represents a user in the system
type User struct {
	BaseModel
	Email           string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password        string     `gorm:"not null;size:255" json:"-"`
	FirstName       string     `gorm:"size:100" json:"first_name"`
	LastName        string     `gorm:"size:100" json:"last_name"`
	Role            UserRole   `gorm:"type:varchar(20);default:user" json:"role"`
	Status          UserStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	RefreshToken    string     `gorm:"size:500" json:"-"`

	// Profile fields
	Avatar      string `gorm:"size:500" json:"avatar,omitempty"`
	Bio         string `gorm:"size:1000" json:"bio,omitempty"`
	PhoneNumber string `gorm:"size:20" json:"phone_number,omitempty"`

	// Relations
	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Hash password if not already hashed
	if len(u.Password) > 0 && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

// BeforeUpdate hook for User
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Check if password is being updated
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEmailVerified checks if the user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// UserResponse is the response structure for user data (without sensitive fields)
type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	FullName        string     `json:"full_name"`
	Role            UserRole   `json:"role"`
	Status          UserStatus `json:"status"`
	Avatar          string     `json:"avatar,omitempty"`
	Bio             string     `json:"bio,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		FullName:        u.FullName(),
		Role:            u.Role,
		Status:          u.Status,
		Avatar:          u.Avatar,
		Bio:             u.Bio,
		PhoneNumber:     u.PhoneNumber,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
