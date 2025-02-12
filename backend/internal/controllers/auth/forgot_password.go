package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type resetPasswordPayload struct {
	token       string
	newPassword string
}

// Endpoint to handle forgot password
func (a *AuthController) ForgotPasswordHandler(c *gin.Context) {
	// 1) Unmarshal request and sanitise email
	var request forgotPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid Request Payload", err)
		return
	}

	// 2) Validate the email
	validator := validator.New()
	if err := validator.Struct(request); err != nil {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid email", err)
		return
	}

	// 2) Check if the user exists
	user, err := a.UserDB.FindUserByEmail(request.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database lookup error", err)
		return
	}

	// 3) If user doesn't exist - tell the user if an account exists with that email an email has been sent
	if user == nil {
		a.HandleError(c, http.StatusUnauthorized, "Invalid email or password", "User not found", nil)
		return
	}

	// 4) If user actually does exist, generate link

	// 5) Send link to users email
	link := "http://wealthscope/forgot-password/12hkj8lk9sdfjkl23jk"
	a.sendForgotPasswordToken(user.Email, link)

	// 6) Return 200 ok
}

// Endpoint to handle resetting password
func (a *AuthController) ResetPasswordHandler(c *gin.Context) {
	// 1) Unmarshal token and new password

	// 2) Check if token exists and it hasn't expired

	// 3) Hash and update the password in the database

	// 4) Return 200 ok
}

// Create forgot password token
func generateUuid() {

}

// Send forgot password token
func sendToken() {

}
