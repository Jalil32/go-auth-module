package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type registerRequest struct {
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,strong_password"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type ValidationError struct {
	UserMessage   string            // General user-friendly message
	FieldErrors   map[string]string // Field-specific validation errors
	InternalError error             // Internal error for logging
}

// Validate validates the requestPayload using the validator
func (rp *registerRequest) Validate() *ValidationError {
	validate := validator.New()
	validate.RegisterValidation("strong_password", passwordValidator)

	err := validate.Struct(rp)
	if err != nil {
		fieldErrors := make(map[string]string)

		// Iterate over validation errors
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field() // Field name (e.g., "Email", "Password")
			tag := err.Tag()     // Validation rule that failed (e.g., "required", "email", "strong_password")

			// Customize the error message based on the field and tag
			switch tag {
			case "required":
				fieldErrors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				fieldErrors[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "strong_password":
				fieldErrors[field] = fmt.Sprintf("%s must be at least 8 characters long and contain at least one lowercase letter, one uppercase letter, one digit, and one special character", field)
			default:
				fieldErrors[field] = fmt.Sprintf("%s failed validation: %s", field, tag)
			}
		}

		return &ValidationError{
			UserMessage:   "Validation failed",
			FieldErrors:   fieldErrors,
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

func (a *AuthController) Register(c *gin.Context) {
	// 1) Unmarshal request body into request struct
	var registerRequest registerRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// 2) Validate the requestPayload
	if validationErr := registerRequest.Validate(); validationErr != nil {
		// Construct a detailed error response
		errorResponse := gin.H{
			"message": validationErr.UserMessage,
			"errors":  validationErr.FieldErrors,
		}

		// Include internal error details in test mode
		if gin.Mode() == gin.TestMode {
			errorResponse["internal_message"] = validationErr.InternalError.Error()
		}

		// Log the error and send the response
		a.HandleError(c, http.StatusBadRequest, validationErr.UserMessage, "Validation failed", validationErr.InternalError)
		return
	}

	// 3) Check if user exists already
	existingUser, err := a.UserDB.FindUserByEmail(registerRequest.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	if existingUser != nil {
		a.HandleError(c, http.StatusConflict, "User with this email already exists.", "User already exists in db", nil)
		return
	}

	// 4) Hash password
	hashedPassword, err := a.hashPassword(registerRequest.Password)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to hash password", err)
		return
	}

	// 5) Create new user
	newUser := models.User{
		Email:        registerRequest.Email,
		FirstName:    registerRequest.FirstName,
		LastName:     registerRequest.LastName,
		PasswordHash: &hashedPassword,
	}

	// 6) Start transaction
	tx, err := a.UserDB.Beginx()
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to start transaction", err)
		return
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.Logger.Error("Failed to rollback transaction", "error", rbErr)
			}
		}
	}()

	// 7) Create user in transaction
	if err := a.UserDB.CreateUser(tx, &newUser); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to create user", err)
		return
	}

	// 8) Send OTP
	if err := a.sendOTP(newUser.Email); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to send OTP", err)
		return
	}

	// 9) Commit transaction
	if err := tx.Commit(); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to commit transaction", err)
		return
	}

	// 10) Send success response
	c.JSON(http.StatusCreated, gin.H{"message": "User created and OTP sent successfully"})
}
