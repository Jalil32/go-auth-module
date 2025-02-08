package auth

import (
	"net/http"
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
	existingUser, err := a.UserDB.FindUserByEmail(req.Email)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Database error during user lookup", err)
		return
	}

	if existingUser != nil {
		a.handleError(c, http.StatusConflict, "User already exists", nil)
		return
	}

	// 3) Hash password
	hashedPassword, err := a.hashPassword(req.Password)
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// 4) Create new user
	newUser := models.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: &hashedPassword,
	}

	// 5) Start transaction
	tx, err := a.UserDB.Beginx()
	if err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to start transaction", err)
		return
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.Logger.Error("Failed to rollback transaction", "error", rbErr)
			}
		}
	}()

	// 6) Create user in transaction
	if err := a.UserDB.CreateUser(tx, &newUser); err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// 7) Send OTP
	if err := a.sendOTP(newUser.Email); err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to send OTP", err)
		return
	}

	// 8) Commit transaction
	if err := tx.Commit(); err != nil {
		a.handleError(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	// 9) Send success response
	c.JSON(http.StatusOK, gin.H{"message": "User created and OTP sent successfully"})
}

func (a *AuthController) handleError(c *gin.Context, statusCode int, message string, err error) {
	a.Logger.Error(message, "error", err)
	c.JSON(statusCode, gin.H{"error": message})
}
