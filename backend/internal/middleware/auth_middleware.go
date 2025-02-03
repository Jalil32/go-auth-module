package middleware

import (
	"fmt"
	"net/http"
	"time"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type user struct {
	// add the struct in
}

// AuthMiddleware is the middleware that checks for the presence and validity of the JWT token
func (m *Middleware) AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from cookie or header
		token, err := c.Cookie("auth_token") // Try getting it from cookies
		if err != nil {
			token = c.GetHeader("Authorization") // Try getting from Authorization header
			if token == "" {
				m.Logger.Error("Missing token", "error", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
				c.Abort()
				return
			}
		}

		// Decode the JWT
		user, err := m.decodeJWT(token, secretKey)
		if err != nil {
			m.Logger.Error("Invalid or expired token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user info in context for later use
		c.Set("user", user)
		c.Next()
	}
}

// decodeJWT extracts user information from the JWT token
func (m *Middleware) decodeJWT(tokenString string, secretKey string) (*models.User, error) {
	// Parse the token with claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			m.Logger.Error("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil // Return the secret key
	})

	// Check for parsing errors
	if err != nil {
		m.Logger.Error("Failed to parse token", "error", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &models.User{}

		// Extract user ID
		if userID, ok := claims["user_id"].(float64); ok {
			user.ID = int(userID) // Convert float64 to int
		} else {
			m.Logger.Error("Invalid user id in token")
			return nil, fmt.Errorf("invalid user_id type in token")
		}

		// Extract email
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		} else {
			m.Logger.Error("Invalid email in token")
			return nil, fmt.Errorf("invalid email type in token")
		}

		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				m.Logger.Error("Token has expired")
				return nil, fmt.Errorf("token has expired")
			}
		}

		return user, nil
	}

	m.Logger.Error("Invalid token")
	return nil, fmt.Errorf("invalid token")
}
