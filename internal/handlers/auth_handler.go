package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-enterprise-api/internal/middleware"
	"github.com/yourusername/go-enterprise-api/internal/services"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/response"
	"github.com/yourusername/go-enterprise-api/pkg/validator"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents the refresh token request body
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate request
	v := validator.New()
	v.Required("email", req.Email, "")
	v.Email("email", req.Email, "")
	v.Required("password", req.Password, "")
	v.Password("password", req.Password)
	v.Required("first_name", req.FirstName, "")
	v.Required("last_name", req.LastName, "")

	if errs := v.Validate(); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	// Call service
	serviceReq := &services.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, tokens, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, gin.H{
		"user":   user.ToResponse(),
		"tokens": tokens,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate request
	v := validator.New()
	v.Required("email", req.Email, "")
	v.Email("email", req.Email, "")
	v.Required("password", req.Password, "")

	if errs := v.Validate(); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	// Call service
	serviceReq := &services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, tokens, err := h.authService.Login(c.Request.Context(), serviceReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"user":   user.ToResponse(),
		"tokens": tokens,
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate user's refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	user := middleware.MustGetUser(c)

	if err := h.authService.Logout(c.Request.Context(), user.ID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out successfully", nil)
}

// RefreshTokens handles token refresh
// @Summary Refresh tokens
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	tokens, err := h.authService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"tokens": tokens,
	})
}

// Me returns the current authenticated user
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	user := middleware.MustGetUser(c)
	response.Success(c, gin.H{
		"user": user.ToResponse(),
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change the current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate new password
	v := validator.New()
	v.Password("new_password", req.NewPassword)

	if errs := v.Validate(); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	user := middleware.MustGetUser(c)

	if err := h.authService.ChangePassword(c.Request.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
		if apperrors.IsAppError(err) {
			response.Error(c, err)
			return
		}
		response.Error(c, apperrors.ErrInternal)
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}
