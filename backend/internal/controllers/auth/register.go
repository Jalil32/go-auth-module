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
	// 1) Unmarshal request body into request struct
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.handleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// 2) Check if user exists already
	existingUser, err := db.FindUserByEmail(a.DB, req.Email)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Database error during user lookup", err)
		return
	}

	if existingUser != nil {
		a.handleError(c, http.StatusConflict, "User already exists", nil)
		return
	}

	// 3) If user does not exist hash password and create new user
	hashedPassword, err := a.hashPassword(req.Password)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	newUser := models.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: &hashedPassword,
	}

	// 4)  Start transaction
	tx, err := a.DB.Beginx()

	if err := db.CreateUser(tx, &newUser); err != nil {
		tx.Rollback()
		a.handleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// 5) Send OTP to new users email, if this fails we rollback the transaction
	err = a.sendOTP(newUser.Email)
	if err != nil {
		tx.Rollback()
		a.handleError(c, http.StatusInternalServerError, "Failed to send otp", err)
		return
	}

	// 6) Commit the user and send back success
	tx.Commit()

	a.Logger.Info("User registered successfully and otp send", "email", newUser.Email, "userID", newUser.ID)
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
