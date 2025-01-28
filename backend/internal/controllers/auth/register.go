package auth

import (
	"net/http"
	"wealthscope/backend/internal/db"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Register handles user registration.
func (a *AuthController) Register(c *gin.Context) {
	var newUser models.User

	// Bind the JSON request body to the newUser struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		a.Logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if the user already exists
	existingUser, err := db.FindUserByEmail(a.DB, newUser.Email)
	if err != nil {
		a.Logger.Error("Database error during user lookup", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during user lookup"})
		return
	}

	if existingUser != nil {
		a.Logger.Error("User already exists", "email", newUser.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	if newUser.PasswordHash == nil {
		a.Logger.Error("No password given")
		c.JSON(http.StatusConflict, gin.H{"error": "Please enter a password"})
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := hashPassword(*newUser.PasswordHash)
	if err != nil {
		a.Logger.Error("Failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser.PasswordHash = &hashedPassword

	// Create the new user in the database
	err = db.CreateUser(a.DB, &newUser)
	if err != nil {
		a.Logger.Error("Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := a.generateJWT(&newUser)
	if err != nil {
		a.Logger.Error("Failed to generate JWT token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the JWT token as a cookie
	a.setAuthCookie(c, token)

	a.Logger.Info("User registered successfully", "email", newUser.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
