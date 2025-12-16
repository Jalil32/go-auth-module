package auth

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	Validate() *ValidationError
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,strong_password"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,strong_password"`
}

type ValidationError struct {
	UserMessage   string            // General user-friendly message
	FieldErrors   map[string]string // Field-specific validation errors
	InternalError error             // Internal error for logging
}

func (rr *LoginRequest) Validate() *ValidationError {
	return validateStruct(rr)
}

func (r *RegisterRequest) Validate() *ValidationError {
	return validateStruct(r)
}

func (fpr *ForgotPasswordRequest) Validate() *ValidationError {
	return validateStruct(fpr)
}

func (fpr *ResetPasswordRequest) Validate() *ValidationError {
	return validateStruct(fpr)
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("strong_password", passwordValidator); err != nil {
		panic(fmt.Sprintf("failed to register password validator: %v", err))
	}
}

func validateStruct(s interface{}) *ValidationError {
	err := validate.Struct(s)
	if err != nil {
		var errorMessages []string

		// Iterate over validation errors
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return &ValidationError{
				UserMessage:   "Validation failed",
				InternalError: err,
			}
		}

		for _, err := range validationErrors {
			field := err.Field() // Field name (e.g., "Email", "Password")
			tag := err.Tag()     // Validation rule that failed (e.g., "required", "email", "strong_password")

			// Customize the error message based on the field and tag
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("The %s field is required.", field)
			case "email":
				message = fmt.Sprintf("The %s field must be a valid email address.", field)
			case "strong_password":
				message = fmt.Sprintf("The %s field must be at least 8 characters long and contain at least one lowercase letter, one uppercase letter, one digit, and one special character.", field)
			case "min":
				message = fmt.Sprintf("The %s field must be at least %s characters long.", field, err.Param())
			case "max":
				message = fmt.Sprintf("The %s field must be no more than %s characters long.", field, err.Param())
			case "eqfield":
				message = fmt.Sprintf("The %s field must match the %s field.", field, err.Param())
			default:
				message = fmt.Sprintf("The %s field is invalid: %s.", field, tag)
			}

			errorMessages = append(errorMessages, message)
		}

		userMessage := "Please correct the following errors:\n"
		for _, msg := range errorMessages {
			userMessage += fmt.Sprintf("- %s\n", msg)
		}

		return &ValidationError{
			UserMessage:   userMessage,
			InternalError: err,
		}
	}

	return nil
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// Check for at least one digit
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return false
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
		return false
	}

	return true
}
