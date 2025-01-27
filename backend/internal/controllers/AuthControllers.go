package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"wealthscope/backend/config"
	"wealthscope/backend/internal/db"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth/gothic"
)

type AuthController struct {
	DB        *sqlx.DB
	Logger    *slog.Logger
	JwtToken  string
	JwtExpiry string
}

// NewAuthController initializes a new AuthController
func NewAuthController(db *sqlx.DB, logger *slog.Logger) (*AuthController, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return nil, fmt.Errorf("failed to load config: %w", err) // Error message clarified
	}

	return &AuthController{
		DB:        db,
		Logger:    logger,
		JwtToken:  cfg.JWT.Token,
		JwtExpiry: cfg.JWT.Expiry,
	}, nil
}

// Register handles user registration.
func (a *AuthController) Register(c *gin.Context) {
	// Example error handling if needed (e.g., database or validation errors)
	a.Logger.Info("User registration initiated")
	c.JSON(http.StatusOK, gin.H{"message": "Register successful!"})
}

// Login handles user login.
func (a *AuthController) Login(c *gin.Context) {
	a.Logger.Info("User login initiated")
	c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
}

// Logout handles user logout.
func (a *AuthController) Logout(c *gin.Context) {
	a.Logger.Info("User logout initiated")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful!"})
}

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

// setAuthCookie sets the authentication token in a secure, HTTP-only cookie.
func (a *AuthController) setAuthCookie(c *gin.Context, token string) {
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
