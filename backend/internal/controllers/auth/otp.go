package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (a *AuthController) VerifyOTPHandler(c *gin.Context) {
	// 1) Get the otp and email from req body
	var otpRequest OTPRequest

	if err := c.ShouldBindJSON(&otpRequest); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid request payload ", err)
		return
	}

	// 2) Verify the otp
	validated := a.validateOTP(otpRequest.Email, otpRequest.OTP)

	// 3) If incorrect, send back error message
	if !validated {
		a.HandleError(c, http.StatusUnauthorized, "Incorrect or expired one time password. Please try again.", "Incorrect or expired one time password.", nil)
		return
	}

	// 4) If correct fetch user
	existingUser, err := a.UserDB.FindUserByEmail(otpRequest.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database error during user lookup", err)
		return
	}

	// 5) Update user to be verified
	tx, err := a.UserDB.Beginx()
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to start transaction", err)
		return
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.Logger.Error("Failed to rollback transaction", "error", err)
			}
		}
	}()

	existingUser.Verified = true

	if err := a.UserDB.UpdateUser(tx, existingUser); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to update user", err)
		return
	}

	token, err := a.JWTGenerator.GenerateJWT(existingUser)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate token", err)
		return
	}

	// 7) Commit the transaction
	if err := tx.Commit(); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to commit transaction", err)
		return
	}

	a.setAuthCookie(c, token)

	a.Logger.Info("User email successfully verified", "email", existingUser.Email, "userID", existingUser.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully verified",
	})
}

func (a *AuthController) generateNewOTPHandler(c *gin.Context) {
}
