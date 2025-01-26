package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

// AuthController acts as a wrapper for all auth-related handlers
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register handles user registration.
func (a *AuthController) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Register successful!"})
}

// Login handles user login.
func (a *AuthController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
}

// Logout handles user logout.
func (a *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful!"})
}
