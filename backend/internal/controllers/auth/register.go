package auth

import (
	"errors"
	"net/http"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
)

func (a *AuthController) Register(c *gin.Context) {
	// 1) Unmarshal request body into request struct
	var registerRequest RegisterRequest
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
		a.HandleError(c, http.StatusConflict, "User with this email already exists.", "User already exists in db", errors.New("User already exists in db"))
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
