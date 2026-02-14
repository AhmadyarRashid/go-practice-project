package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Details    string      `json:"details,omitempty"`
	Err        error       `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(statusCode int, code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithData adds additional data to the error
func (e *AppError) WithData(data interface{}) *AppError {
	e.Data = data
	return e
}

// Error codes
const (
	// General errors (1000-1999)
	CodeInternalError    = 1000
	CodeValidationError  = 1001
	CodeNotFound         = 1002
	CodeBadRequest       = 1003
	CodeConflict         = 1004
	CodeTooManyRequests  = 1005

	// Authentication errors (2000-2999)
	CodeUnauthorized     = 2000
	CodeInvalidToken     = 2001
	CodeTokenExpired     = 2002
	CodeInvalidCredentials = 2003
	CodeForbidden        = 2004

	// User errors (3000-3999)
	CodeUserNotFound     = 3000
	CodeUserExists       = 3001
	CodeInvalidPassword  = 3002
	CodeEmailExists      = 3003

	// Database errors (4000-4999)
	CodeDatabaseError    = 4000
	CodeRecordNotFound   = 4001
	CodeDuplicateEntry   = 4002
)

// Predefined errors
var (
	// General errors
	ErrInternal = NewAppError(http.StatusInternalServerError, CodeInternalError, "Internal server error")
	ErrNotFound = NewAppError(http.StatusNotFound, CodeNotFound, "Resource not found")
	ErrBadRequest = NewAppError(http.StatusBadRequest, CodeBadRequest, "Bad request")
	ErrValidation = NewAppError(http.StatusBadRequest, CodeValidationError, "Validation error")
	ErrConflict = NewAppError(http.StatusConflict, CodeConflict, "Resource conflict")
	ErrTooManyRequests = NewAppError(http.StatusTooManyRequests, CodeTooManyRequests, "Too many requests")

	// Authentication errors
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, CodeUnauthorized, "Unauthorized")
	ErrInvalidToken = NewAppError(http.StatusUnauthorized, CodeInvalidToken, "Invalid token")
	ErrTokenExpired = NewAppError(http.StatusUnauthorized, CodeTokenExpired, "Token expired")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, CodeInvalidCredentials, "Invalid credentials")
	ErrForbidden = NewAppError(http.StatusForbidden, CodeForbidden, "Forbidden")

	// User errors
	ErrUserNotFound = NewAppError(http.StatusNotFound, CodeUserNotFound, "User not found")
	ErrUserExists = NewAppError(http.StatusConflict, CodeUserExists, "User already exists")
	ErrEmailExists = NewAppError(http.StatusConflict, CodeEmailExists, "Email already exists")
	ErrInvalidPassword = NewAppError(http.StatusBadRequest, CodeInvalidPassword, "Invalid password")

	// Database errors
	ErrDatabase = NewAppError(http.StatusInternalServerError, CodeDatabaseError, "Database error")
	ErrRecordNotFound = NewAppError(http.StatusNotFound, CodeRecordNotFound, "Record not found")
	ErrDuplicateEntry = NewAppError(http.StatusConflict, CodeDuplicateEntry, "Duplicate entry")
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError converts an error to AppError
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return ErrInternal.WithError(err)
}

// ValidationErrors holds multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// Add adds a validation error
func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// Error implements the error interface
func (v *ValidationErrors) Error() string {
	return "validation errors"
}
