package routes

import (
	"log/slog"
	"github.com/jalil32/go-auth-module/config"
	"github.com/jalil32/go-auth-module/internal/controllers/auth"
	"github.com/jalil32/go-auth-module/internal/controllers/bank"
	"github.com/jalil32/go-auth-module/internal/controllers/stock"
	"github.com/jalil32/go-auth-module/internal/db"
	"github.com/jalil32/go-auth-module/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Routes(router *gin.Engine, database *sqlx.DB, rdb *redis.Client, logger *slog.Logger, cfg *config.Config) error {
	// Create user database
	userDB := &db.UserDB{DB: database}
	jwtService := &auth.JWTService{SecretKey: cfg.JWT.Token, JwtExpiry: cfg.JWT.Expiry}

	// Initialise Auth Controller instance
	authController, err := auth.NewAuthController(userDB, rdb, logger, jwtService, cfg)

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
			auth.POST("/forgot-password", authController.ForgotPasswordHandler)
			auth.POST("/reset-password", authController.ResetPasswordHandler)
		}

		stock := api.Group("/stock")
		{
			stock.GET(":symbol", stockController.GetStockQuoteHandler)
		}

		bank := api.Group("/bank")
		{
			bank.POST("/upload", bankController.UploadBankStatement)
		}

		// test endpoint, remove after use
		api.GET("/test", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "hello from backend test endpoint",
			})
		})
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
