package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/config"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/repository"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
	TokenType string         `json:"token_type"`
	jwt.RegisteredClaims
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthService interface defines authentication methods
type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*models.User, *TokenPair, error)
	Login(ctx context.Context, req *LoginRequest) (*models.User, *TokenPair, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetUserFromToken(ctx context.Context, claims *Claims) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
}

// authService implements AuthService
type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*models.User, *TokenPair, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("Failed to check if user exists", logger.Err(err))
		return nil, nil, apperrors.ErrInternal
	}
	if exists {
		return nil, nil, apperrors.ErrEmailExists
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
		Status:    models.StatusActive, // In production, this might be StatusPending until email is verified
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user", logger.Err(err))
		return nil, nil, apperrors.ErrInternal
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		logger.Error("Failed to generate tokens", logger.Err(err))
		return nil, nil, apperrors.ErrInternal
	}

	// Save refresh token
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		logger.Error("Failed to save refresh token", logger.Err(err))
		return nil, nil, apperrors.ErrInternal
	}

	logger.Info("User registered successfully", logger.String("email", user.Email))
	return user, tokens, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req *LoginRequest) (*models.User, *TokenPair, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, apperrors.ErrInvalidCredentials
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, nil, apperrors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, nil, apperrors.ErrForbidden.WithDetails("Account is not active")
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		logger.Error("Failed to generate tokens", logger.Err(err))
		return nil, nil, apperrors.ErrInternal
	}

	// Save refresh token and update last login
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		logger.Error("Failed to save refresh token", logger.Err(err))
	}
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		logger.Error("Failed to update last login", logger.Err(err))
	}

	logger.Info("User logged in successfully", logger.String("email", user.Email))
	return user, tokens, nil
}

// Logout logs out a user
func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	// Clear refresh token
	if err := s.userRepo.UpdateRefreshToken(ctx, userID, ""); err != nil {
		logger.Error("Failed to clear refresh token", logger.Err(err))
		return apperrors.ErrInternal
	}
	return nil
}

// RefreshTokens refreshes the access token using a refresh token
func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	// Check token type
	if claims.TokenType != "refresh" {
		return nil, apperrors.ErrInvalidToken
	}

	// Find user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Verify refresh token matches stored token
	if user.RefreshToken != refreshToken {
		return nil, apperrors.ErrInvalidToken
	}

	// Check if user is still active
	if !user.IsActive() {
		return nil, apperrors.ErrForbidden.WithDetails("Account is not active")
	}

	// Generate new tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	// Save new refresh token
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		logger.Error("Failed to save refresh token", logger.Err(err))
		return nil, apperrors.ErrInternal
	}

	return tokens, nil
}

// ValidateToken validates a JWT token
func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.ErrInvalidToken
}

// GetUserFromToken retrieves user from token claims
func (s *authService) GetUserFromToken(ctx context.Context, claims *Claims) (*models.User, error) {
	return s.userRepo.FindByID(ctx, claims.UserID)
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return apperrors.ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrInternal
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return apperrors.ErrInternal
	}

	// Invalidate all sessions by clearing refresh token
	if err := s.userRepo.UpdateRefreshToken(ctx, userID, ""); err != nil {
		logger.Error("Failed to clear refresh token", logger.Err(err))
	}

	return nil
}

// generateTokenPair generates access and refresh tokens
func (s *authService) generateTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(time.Duration(s.config.JWT.ExpiryHours) * time.Hour)
	refreshExpiry := now.Add(time.Duration(s.config.JWT.RefreshExpiryHours) * time.Hour)

	// Access token claims
	accessClaims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.App.Name,
			Subject:   user.ID.String(),
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh token claims
	refreshClaims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.App.Name,
			Subject:   user.ID.String(),
		},
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}
