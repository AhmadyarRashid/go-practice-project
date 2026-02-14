package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
	"go.uber.org/zap"
)

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestLogger creates a request logging middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Get path and method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		}

		// Add query string if exists
		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}

		// Add user ID if authenticated
		if user, exists := GetUser(c); exists {
			fields = append(fields, zap.String("user_id", user.ID.String()))
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Error("Server error", fields...)
		case statusCode >= 400:
			logger.Warn("Client error", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}
