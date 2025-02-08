package routes

import (
	"log/slog"
	"wealthscope/backend/internal/controllers/auth"
	"wealthscope/backend/internal/controllers/bank"
	"wealthscope/backend/internal/controllers/stock"
	"wealthscope/backend/internal/db"
	"wealthscope/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Routes(router *gin.Engine, database *sqlx.DB, rdb *redis.Client, logger *slog.Logger) error {
	// Create user database
	userDB := &db.UserDB{DB: database}

	// Initialise Auth Controller instance
	authController, err := auth.NewAuthController(userDB, rdb, logger)

	if err != nil {
		logger.Error("Failed to initialise AuthController", "error", err)
		return err
	}

	middleware := middleware.NewMiddlewareSetup(logger)

	// Initialise Stock Controller instance
	stockController := stock.NewStockController(logger)

	// Initialise Bank Controller instance
	bankController := bank.NewBankController(logger, database)

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
			auth.POST("/verify", authController.VerifyOTPHandler)
		}

		stock := api.Group("/stock")
		{
			stock.GET(":symbol", stockController.GetStockQuoteHandler)
		}

		bank := api.Group("/bank")
		{
			bank.POST("/upload", bankController.UploadBankStatement)
		}

	}

	router.GET("/protected", middleware.AuthMiddleware(authController.JwtToken), func(c *gin.Context) {
		// Protected route logic
		user, _ := c.Get("user")
		c.JSON(200, gin.H{
			"user": user,
		})
	})

	return nil
}
