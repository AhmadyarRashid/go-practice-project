package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/response"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]*clientInfo
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type clientInfo struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientInfo),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, info := range rl.requests {
			if now.After(info.resetTime) {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	info, exists := rl.requests[key]
	if !exists || now.After(info.resetTime) {
		rl.requests[key] = &clientInfo{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if info.count >= rl.limit {
		return false
	}

	info.count++
	return true
}

// RateLimit creates a rate limiting middleware
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// Use client IP as key (can be customized to use user ID for authenticated users)
		key := c.ClientIP()

		// Check if authenticated user
		if user, exists := GetUser(c); exists {
			key = user.ID.String()
		}

		if !limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    apperrors.CodeTooManyRequests,
					Message: "Rate limit exceeded. Please try again later.",
				},
			})
			return
		}

		c.Next()
	}
}

// StrictRateLimit creates a stricter rate limiting middleware for sensitive endpoints
func StrictRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// Use combination of IP and path for more granular limiting
		key := c.ClientIP() + ":" + c.Request.URL.Path

		if !limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    apperrors.CodeTooManyRequests,
					Message: "Too many requests to this endpoint. Please try again later.",
				},
			})
			return
		}

		c.Next()
	}
}
