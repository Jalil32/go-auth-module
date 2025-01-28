package auth

import (
	"fmt"
	"net/http"
	"wealthscope/backend/internal/db"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// SignInWithProvider handles third-party sign-in using a provider (e.g., Google)
func (a *AuthController) SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		a.Logger.Error("Provider not specified", "error", "Provider missing in URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not specified"})
		return
	}

	// Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Begin OAuth flow
	a.Logger.Info("Starting OAuth flow", "provider", provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler handles the OAuth callback and user creation
func (a *AuthController) CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		a.Logger.Error("Provider not specified", "error", "Provider missing in URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not specified"})
		return
	}

	// Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		a.Logger.Error("OAuth complete error", "provider", provider, "error", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OAuth complete error: %w", err))
		return
	}

	// Check if the user exists in the database
	existingUser, err := db.FindUserByEmail(a.DB, oauthUser.Email)
	if err != nil {
		a.Logger.Error("Database error during user lookup", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error during user lookup"})
		return
	}

	// Create a new user if not found
	var user *models.User
	if existingUser == nil {
		newUser := models.User{
			Email:     oauthUser.Email,
			FirstName: oauthUser.FirstName,
			LastName:  oauthUser.LastName,
			Provider:  &oauthUser.Provider,
		}
		err := db.CreateUser(a.DB, &newUser)
		if err != nil {
			a.Logger.Error("Failed to create user", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %v", err)})
			return
		}
		user = &newUser
	} else {
		user = existingUser
	}

	// Generate JWT token
	token, err := a.generateJWT(user)
	if err != nil {
		a.Logger.Error("Failed to generate JWT token", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate token: %v", err)})
		return
	}

	// Set the JWT token as a cookie
	a.setAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
