package auth

import (
	"net/http"
	"wealthscope/backend/internal/db"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (a *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.handleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	existingUser, err := db.FindUserByEmail(a.DB, req.Email)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Database error during user lookup", err)
		return
	}

	if existingUser != nil {
		a.handleError(c, http.StatusConflict, "User already exists", nil)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	newUser := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName, PasswordHash: &hashedPassword}
	if err := db.CreateUser(a.DB, &newUser); err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	token, err := a.generateJWT(&newUser)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	a.setAuthCookie(c, token)

	a.Logger.Info("User registered successfully", "email", newUser.Email, "userID", newUser.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"email":     newUser.Email,
			"firstName": newUser.FirstName,
			"lastName":  newUser.LastName,
		},
	})
}

func (a *AuthController) handleError(c *gin.Context, statusCode int, message string, err error) {
	a.Logger.Error(message, "error", err)
	c.JSON(statusCode, gin.H{"error": message})
}
