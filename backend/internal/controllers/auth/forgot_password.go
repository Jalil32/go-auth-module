package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	Validate() *ValidationError
}

type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type resetPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,strong_password"`
}

// forgotPasswordRequest implements the Validatable interface
func (rp *forgotPasswordRequest) Validate() *ValidationError {
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

// resetPasswordRequest implements the Validatable interface
func (rp *resetPasswordRequest) Validate() *ValidationError {
	validate := validator.New()
	validate.RegisterValidation("strong_password", passwordValidator)

	err := validate.Struct(rp)
	if err != nil {
		fieldErrors := make(map[string]string)

		// Iterate over validation errors
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field() // Field name (e.g., "NewPassword")
			tag := err.Tag()     // Validation rule that failed (e.g., "required", "min")

			// Customize the error message based on the field and tag
			switch tag {
			case "required":
				fieldErrors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				fieldErrors[field] = fmt.Sprintf("%s must be at least 8 characters long", field)
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

// Endpoint to handle forgot password
func (a *AuthController) ForgotPasswordHandler(c *gin.Context) {
	// 1) Unmarshal request
	var request forgotPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// 2) Validate the request
	if validationErr := request.Validate(); validationErr != nil {
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

	// 3) Check if the user exists
	user, err := a.UserDB.FindUserByEmail(request.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	if user == nil {
		a.HandleError(c, http.StatusUnauthorized, "If your account exists, we have sent you an email to reset your password.", "User not found", errors.New("No User"))
		return
	}

	// 4) Check if they use a provider and not traditional login/signup
	if user.Provider != nil {
		a.HandleError(c, http.StatusBadRequest, "Please sign in with Google.", "User not found", nil)
		return
	}

	// 5) Generate forgot password link
	link, err := a.generateForgotPasswordLink(request.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate reset link", err)
		return
	}

	// 6) Send the forgot password link to the user
	if err := a.sendForgotPasswordToken(request.Email, link); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to send reset email", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If your account exists, we have sent you an email to reset your password.",
	})
}

func (a *AuthController) ResetPasswordHandler(c *gin.Context) {
	// 1) Extract token from URL parameters
	token := c.Query("token")
	if token == "" {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Token is required", errors.New("Missing token"))
		return
	}
	// 1) Unmarshal the request
	var request resetPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// 2) Validate the request, we need to check that the new password fulfills our requirements
	if validationErr := request.Validate(); validationErr != nil {
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

	// 4) Validate the forgot password token
	email, err := a.validateForgotPasswordToken(token)
	if err != nil {
		a.HandleError(c, http.StatusBadRequest, "Invalid or expired token", "Token validation failed", err)
		return
	}

	// 5) Hash the new password
	hashedPassword, err := a.hashPassword(request.NewPassword)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to hash password", err)
		return
	}

	// 6) Get the users details
	user, err := a.UserDB.FindUserByEmail(email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	// 7) Set the users password to new password hash
	user.PasswordHash = &hashedPassword

	// 8) Update the user
	tx, err := a.UserDB.Beginx()
	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.Logger.Error("Failed to rollback transaction", "error", rbErr)
			}
		}
	}()

	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to start transaction", err)
		return
	}

	if err := a.UserDB.UpdateUser(tx, user); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to update password", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to commit transaction", err)
		return
	}

	// 9) Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully.",
	})
}
