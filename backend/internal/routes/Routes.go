package routes

import (
	"log/slog"
	"wealthscope/backend/internal/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Routes(router *gin.Engine, db *sqlx.DB, logger *slog.Logger) error {
	// Initialise Auth Controller instance
	authController, err := controllers.NewAuthController(db, logger)

	if err != nil {
		logger.Error("Failed to initialise AuthController", "error", err)
		return err
	}

	// Register controllers to routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
			auth.GET("/:provider", authController.SignInWithProvider)
			auth.GET("/:provider/callback", authController.CallbackHandler)
		}
	}

	return nil
}
