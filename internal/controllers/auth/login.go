package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (a *AuthController) Login(c *gin.Context) {
	var loginRequest LoginRequest

	// 1) Bind and validate the request
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// 2) Validate the request
	if validationErr := loginRequest.Validate(); validationErr != nil {
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

	// 3) Find the user by email
	user, err := a.UserDB.FindUserByEmail(loginRequest.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	if user == nil {
		// No user found with the given email
		a.HandleError(c, http.StatusUnauthorized, "Invalid email or password", "User not found", errors.New("User not found"))
		return
	}

	// 4) Check if the user needs to authenticate via a provider
	if user.PasswordHash == nil && user.Provider != nil {
		a.HandleError(c, http.StatusUnauthorized, "Please sign in with a provider", "User needs to sign in with provider", errors.New("User needs to sign in with provider"))
		return
	}

	// 5) Compare the provided password with the hashed password
	if compareErr := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(loginRequest.Password)); compareErr != nil {
		a.HandleError(c, http.StatusUnauthorized, "Invalid email or password", "Invalid password", compareErr)
		return
	}

	// 6) Handle unverified users
	if !user.Verified {
		if otpErr := a.sendOTP(user.Email); otpErr != nil {
			a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to send OTP", otpErr)
			return
		}
		// Not using normal error handling as this is a special case
		a.Logger.Info("User not verified. OTP has been sent.", "email", user.Email, "userID", user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not verified.",
			"user": gin.H{
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
			},
		})
		return
	}

	// 7) Generate and set JWT token
	token, err := a.JWTGenerator.GenerateJWT(user)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate JWT Token", err)
		return
	}

	a.setAuthCookie(c, token)

	a.Logger.Info("User logged in successfully", "email", loginRequest.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
