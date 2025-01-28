package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout handles user logout.
func (a *AuthController) Logout(c *gin.Context) {
	// Clear the authentication cookie
	c.SetCookie(
		"auth_token", // Name
		"",           // Value
		-1,           // MaxAge (expire immediately)
		"/",          // Path
		"",           // Domain (empty for default)
		true,         // Secure (true for HTTPS)
		true,         // HttpOnly
	)

	a.Logger.Info("User logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
