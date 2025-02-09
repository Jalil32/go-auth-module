package auth

import "github.com/gin-gonic/gin"

func (a *AuthController) HandleError(c *gin.Context, statusCode int, userMessage string, internalMessage string, err error) {
	a.Logger.Error("internal message", internalMessage, "user message", userMessage, "error", err)
	c.JSON(statusCode, gin.H{"message": userMessage})
}
