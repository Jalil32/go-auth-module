package auth

import (
	"fmt"
	"net/http"
	"time"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// generateJWT creates a JWT token for the authenticated user
func (a *AuthController) generateJWT(user *models.User) (string, error) {
	// Parse the expiry time from the string (e.g., "10m")
	expiryTime, err := time.ParseDuration(a.JwtExpiry)
	if err != nil {
		a.Logger.Error("Failed to parse expiry time", "error", err)
		return "", fmt.Errorf("failed to parse expiry time: %w", err)
	}

	// Calculate expiry time in Unix seconds
	expiryUnix := time.Now().Add(expiryTime).Unix()

	// Create the JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expiryUnix, // Use the actual Unix timestamp for expiry
		"iat":     time.Now().Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString([]byte(a.JwtToken))
	if err != nil {
		a.Logger.Error("Failed to sign JWT token", "error", err)
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil // Return the signed token and nil error
}

// setAuthCookie sets the authentication token in a secure, HTTP-only cookie.
func (a *AuthController) setAuthCookie(c *gin.Context, token string) {
	// Clear the _gothic_session cookie as we are not using this
	c.SetCookie("_gothic_session", "", -1, "/", "", true, true)

	// Parse the expiry duration (e.g., "10m")
	expiryDuration, err := time.ParseDuration(a.JwtExpiry)
	if err != nil {
		a.Logger.Error("Invalid cookie expiry duration", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid cookie expiry duration"})
		return
	}

	// Calculate the cookie expiry time in seconds
	expirySeconds := int(expiryDuration.Seconds())

	// Set the HTTP-only cookie
	c.SetCookie(
		"auth_token",  // Name
		token,         // Value
		expirySeconds, // MaxAge in seconds
		"/",           // Path
		"",            // Domain (empty for default)
		true,          // Secure (true for HTTPS)
		true,          // HttpOnly
	)
}
