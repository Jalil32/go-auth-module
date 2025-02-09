package auth

import (
	"fmt"
	"net/http"
	"time"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService implements JWTGenerator
type JWTService struct {
	SecretKey string
}

// GenerateJWT creates a JWT token for the authenticated user
func (j *JWTService) GenerateJWT(user *models.User) (string, error) {
	expiryTime, err := time.ParseDuration("1h") // or use config value
	if err != nil {
		return "", fmt.Errorf("failed to parse expiry time: %w", err)
	}

	expiryUnix := time.Now().Add(expiryTime).Unix()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expiryUnix,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
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
