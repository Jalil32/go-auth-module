package routes

import (
	"wealthscope/backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	// Initialise Auth Controller instance
	authController := controllers.NewAuthController()

	// Register controllers to routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
		}
	}
}
