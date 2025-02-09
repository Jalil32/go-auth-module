package auth

import (
	"net/http"
	"regexp"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type registerRequest struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// Validate validates the requestPayload using the validator
func (rp *registerRequest) validate() error {
	validate := validator.New()
	validate.RegisterValidation("strong_password", passwordValidator)
	return validate.Struct(rp)
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// Check for at least one digit
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return false
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
		return false
	}

	return true
}

func (a *AuthController) Register(c *gin.Context) {
	// 1) Unmarshal request body into request struct
	var registerRequest registerRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// Validate the requestPayload
	if err := registerRequest.validate(); err != nil {
		a.HandleError(c, http.StatusBadRequest, err.Error(), err.Error(), err)
		return
	}

	// 2) Check if user exists already
	existingUser, err := a.UserDB.FindUserByEmail(registerRequest.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	if existingUser != nil {
		a.HandleError(c, http.StatusConflict, "User with this email already exists.", "User already exists in db", nil)
		return
	}

	// 3) Hash password
	hashedPassword, err := a.hashPassword(registerRequest.Password)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to hash password", err)
		return
	}

	// 4) Create new user
	newUser := models.User{
		Email:        registerRequest.Email,
		FirstName:    registerRequest.FirstName,
		LastName:     registerRequest.LastName,
		PasswordHash: &hashedPassword,
	}

	// 5) Start transaction
	tx, err := a.UserDB.Beginx()
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to start transaction", err)
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
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to create user", err)
		return
	}

	// 7) Send OTP
	if err := a.sendOTP(newUser.Email); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to send OTP", err)
		return
	}

	// 8) Commit transaction
	if err := tx.Commit(); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to commit transaction", err)
		return
	}

	// 9) Send success response
	c.JSON(http.StatusOK, gin.H{"message": "User created and OTP sent successfully"})
}
