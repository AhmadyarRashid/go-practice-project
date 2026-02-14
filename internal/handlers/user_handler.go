package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/internal/middleware"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/services"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/response"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAll returns all users with pagination
// @Summary Get all users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert to response
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response.Paginated(c, userResponses, page, pageSize, total)
}

// GetByID returns a user by ID
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"user": user.ToResponse(),
	})
}

// Update updates a user
// @Summary Update user
// @Description Update a user's profile (own profile or admin)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body services.UpdateUserRequest true "Update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// Check if user is updating their own profile or is admin
	currentUser := middleware.MustGetUser(c)
	if currentUser.ID != id && !currentUser.IsAdmin() {
		response.Forbidden(c, "You can only update your own profile")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"user": user.ToResponse(),
	})
}

// Delete deletes a user
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// Prevent admin from deleting themselves
	currentUser := middleware.MustGetUser(c)
	if currentUser.ID == id {
		response.BadRequest(c, "You cannot delete your own account")
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// Search searches for users
// @Summary Search users
// @Description Search for users by name or email
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/search [get]
func (h *UserHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.Search(c.Request.Context(), query, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert to response
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response.Paginated(c, userResponses, page, pageSize, total)
}

// UpdateStatusRequest represents the status update request
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus updates a user's status (admin only)
// @Summary Update user status
// @Description Update a user's account status (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateStatusRequest true "Status data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/status [patch]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate status
	validStatuses := []string{string(models.StatusActive), string(models.StatusInactive), string(models.StatusBanned), string(models.StatusPending)}
	isValid := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValid = true
			break
		}
	}
	if !isValid {
		response.BadRequest(c, "Invalid status value")
		return
	}

	if err := h.userService.UpdateStatus(c.Request.Context(), id, models.UserStatus(req.Status)); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Status updated successfully", nil)
}

// UpdateRoleRequest represents the role update request
type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateRole updates a user's role (admin only)
// @Summary Update user role
// @Description Update a user's role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateRoleRequest true "Role data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/role [patch]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate role
	validRoles := []string{string(models.RoleUser), string(models.RoleAdmin), string(models.RoleModerator)}
	isValid := false
	for _, r := range validRoles {
		if req.Role == r {
			isValid = true
			break
		}
	}
	if !isValid {
		response.BadRequest(c, "Invalid role value")
		return
	}

	// Prevent admin from demoting themselves
	currentUser := middleware.MustGetUser(c)
	if currentUser.ID == id && req.Role != string(models.RoleAdmin) {
		response.Error(c, apperrors.ErrBadRequest.WithDetails("You cannot change your own role"))
		return
	}

	if err := h.userService.UpdateRole(c.Request.Context(), id, models.UserRole(req.Role)); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Role updated successfully", nil)
}
