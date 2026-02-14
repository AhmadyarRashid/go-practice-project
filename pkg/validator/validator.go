package validator

import (
	"regexp"
	"strings"
	"unicode"

	apperrors "github.com/yourusername/go-enterprise-api/pkg/errors"
)

// Validator provides validation utilities
type Validator struct {
	errors *apperrors.ValidationErrors
}

// New creates a new Validator
func New() *Validator {
	return &Validator{
		errors: apperrors.NewValidationErrors(),
	}
}

// Validate returns validation errors if any
func (v *Validator) Validate() *apperrors.ValidationErrors {
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return v.errors.HasErrors()
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors.Add(field, message)
}

// Required validates that a field is not empty
func (v *Validator) Required(field, value, message string) *Validator {
	if strings.TrimSpace(value) == "" {
		if message == "" {
			message = field + " is required"
		}
		v.errors.Add(field, message)
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int, message string) *Validator {
	if len(value) < min {
		if message == "" {
			message = field + " must be at least " + string(rune(min)) + " characters"
		}
		v.errors.Add(field, message)
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int, message string) *Validator {
	if len(value) > max {
		if message == "" {
			message = field + " must be at most " + string(rune(max)) + " characters"
		}
		v.errors.Add(field, message)
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value, message string) *Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		if message == "" {
			message = "Invalid email format"
		}
		v.errors.Add(field, message)
	}
	return v
}

// Password validates password strength
func (v *Validator) Password(field, value string) *Validator {
	var (
		hasMinLen  = len(value) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLen {
		v.errors.Add(field, "Password must be at least 8 characters")
	}
	if !hasUpper {
		v.errors.Add(field, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		v.errors.Add(field, "Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		v.errors.Add(field, "Password must contain at least one number")
	}
	if !hasSpecial {
		v.errors.Add(field, "Password must contain at least one special character")
	}

	return v
}

// Match validates that two values match
func (v *Validator) Match(field1, value1, field2, value2, message string) *Validator {
	if value1 != value2 {
		if message == "" {
			message = field1 + " and " + field2 + " do not match"
		}
		v.errors.Add(field1, message)
	}
	return v
}

// UUID validates UUID format
func (v *Validator) UUID(field, value, message string) *Validator {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(value) {
		if message == "" {
			message = "Invalid UUID format"
		}
		v.errors.Add(field, message)
	}
	return v
}

// InSlice validates that a value is in a slice
func (v *Validator) InSlice(field, value string, allowed []string, message string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	if message == "" {
		message = field + " must be one of: " + strings.Join(allowed, ", ")
	}
	v.errors.Add(field, message)
	return v
}

// Custom adds a custom validation
func (v *Validator) Custom(field string, valid bool, message string) *Validator {
	if !valid {
		v.errors.Add(field, message)
	}
	return v
}

// ValidateEmail validates email format and returns bool
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength and returns bool
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
