package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error information in response
type ErrorInfo struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details string      `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Meta represents pagination or additional metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage sends a success response with a message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Resource created successfully",
		Data:    data,
	})
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, page, perPage int, total int64) {
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	appErr := apperrors.GetAppError(err)

	c.JSON(appErr.StatusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
			Data:    appErr.Data,
		},
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errors *apperrors.ValidationErrors) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeValidationError,
			Message: "Validation failed",
			Data:    errors.Errors,
		},
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeBadRequest,
			Message: message,
		},
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeUnauthorized,
			Message: message,
		},
	})
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeForbidden,
			Message: message,
		},
	})
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeNotFound,
			Message: message,
		},
	})
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    apperrors.CodeInternalError,
			Message: message,
		},
	})
}
