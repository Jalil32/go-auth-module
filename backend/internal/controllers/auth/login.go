package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type requestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user login.
func (a *AuthController) Login(c *gin.Context) {
	loginRequest := &requestPayload{
		Email:    "",
		Password: "",
	}

	// Bind the JSON request body to the loginRequest struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.Logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Find the user by email
	user, err := a.UserDB.FindUserByEmail(loginRequest.Email)
	if err != nil {
		a.Logger.Error("Database error during user lookup", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during user lookup"})
		return
	}

	if user == nil {
		a.Logger.Error("User not found", "email", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.PasswordHash == nil && user.Provider != nil {
		a.Logger.Error("User needs to sign in with provider", "email", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please sign in with a provider"})
		return

	}

	// Compare the provided password with the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		a.Logger.Error("Invalid password", "email", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 5) Send OTP to new users email, if this fails we rollback the transaction
	if !user.Verified {
		err = a.sendOTP(user.Email)
		if err != nil {
			a.handleError(c, http.StatusInternalServerError, "Failed to send otp", err)
			return
		}

		a.Logger.Info("User registered successfully and otp send", "email", user.Email, "userID", user.ID)
		c.JSON(http.StatusOK, gin.H{
			"message": "User is not verified.",
			"user": gin.H{
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
			},
		})
		return
	}

	// Generate JWT token
	token, err := a.generateJWT(user)
	if err != nil {
		a.Logger.Error("Failed to generate JWT token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the JWT token as a cookie
	a.setAuthCookie(c, token)

	a.Logger.Info("User logged in successfully", "email", loginRequest.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
