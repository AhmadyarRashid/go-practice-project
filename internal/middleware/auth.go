package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/services"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/response"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// UserKey is the context key for storing user
	UserKey = "user"
	// ClaimsKey is the context key for storing claims
	ClaimsKey = "claims"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check for Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			response.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// Check if it's an access token
		if claims.TokenType != "access" {
			response.Error(c, apperrors.ErrInvalidToken.WithDetails("Invalid token type"))
			c.Abort()
			return
		}

		// Get user from token
		user, err := authService.GetUserFromToken(c.Request.Context(), claims)
		if err != nil {
			response.Error(c, apperrors.ErrUserNotFound)
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive() {
			response.Forbidden(c, "Account is not active")
			c.Abort()
			return
		}

		// Store user and claims in context
		c.Set(UserKey, user)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It tries to authenticate but doesn't fail if no token is provided
func OptionalAuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			c.Next()
			return
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		if claims.TokenType != "access" {
			c.Next()
			return
		}

		user, err := authService.GetUserFromToken(c.Request.Context(), claims)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserKey, user)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetUser(c)
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireAdminOrModerator creates a middleware that requires admin or moderator role
func RequireAdminOrModerator() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin, models.RoleModerator)
}

// GetUser retrieves the authenticated user from context
func GetUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(UserKey)
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetClaims retrieves the JWT claims from context
func GetClaims(c *gin.Context) (*services.Claims, bool) {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		return nil, false
	}
	return claims.(*services.Claims), true
}

// MustGetUser retrieves the authenticated user from context or panics
func MustGetUser(c *gin.Context) *models.User {
	user, exists := GetUser(c)
	if !exists {
		panic("user not found in context - ensure AuthMiddleware is used")
	}
	return user
}
