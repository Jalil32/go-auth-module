package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func (a *AuthController) HandleError(c *gin.Context, statusCode int, userMessage string, internalMessage string, err error) {
	a.Logger.Error("internal message", internalMessage, "user message", userMessage, "error", err)

	if err == nil {
		err = errors.New("No error")
	}

	// 1) Prepare the response
	response := gin.H{"message": userMessage}

	// 2) If in test mode, include the internal message for debugging
	if gin.Mode() == gin.TestMode {
		response["internal_message"] = internalMessage
		response["error"] = err.Error()
	}

	c.JSON(statusCode, response)
}
