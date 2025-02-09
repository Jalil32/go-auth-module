package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (lr *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(lr)
}

func (a *AuthController) Login(c *gin.Context) {
	var loginRequest LoginRequest

	// 1) Bind and validate the request
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	if err := loginRequest.Validate(); err != nil {
		a.HandleError(c, http.StatusBadRequest, err.Error(), err.Error(), err)
		return
	}

	// 2) Find the user by email
	user, err := a.UserDB.FindUserByEmail(loginRequest.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	if user == nil {
		a.HandleError(c, http.StatusUnauthorized, "Invalid email or password", "Database lookup error", err)
		return
	}

	// 3) Check if the user needs to authenticate via a provider
	if user.PasswordHash == nil && user.Provider != nil {
		a.HandleError(c, http.StatusUnauthorized, "Please sign in with a provider", "User needs to sign in with provider", err)
		return
	}

	// 4) Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		a.HandleError(c, http.StatusUnauthorized, "Invalid email or password", "Invalid password", err)
		return
	}

	// 5) Handle unverified users
	if !user.Verified {
		if err := a.sendOTP(user.Email); err != nil {
			a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to send OTP", err)
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

	// 6) Generate and set JWT token
	token, err := a.JWTGenerator.GenerateJWT(user)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate JWT Token", err)
		return
	}

	a.setAuthCookie(c, token)

	a.Logger.Info("User logged in successfully", "email", loginRequest.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
