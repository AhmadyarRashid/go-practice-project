package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
	"github.com/yourusername/go-enterprise-api/pkg/response"
)

// Recovery creates a recovery middleware that handles panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the panic
				logger.Error("Panic recovered",
					logger.String("error", fmt.Sprintf("%v", err)),
					logger.String("stack", stack),
					logger.String("request_id", GetRequestID(c)),
					logger.String("path", c.Request.URL.Path),
					logger.String("method", c.Request.Method),
				)

				// Respond with error
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Success: false,
					Error: &response.ErrorInfo{
						Code:    apperrors.CodeInternalError,
						Message: "Internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}
