package server

import (
	"log/slog"
	"time"
	"wealthscope/backend/config"
	"wealthscope/backend/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func StartServer(cfg *config.Config, db *sqlx.DB, logger *slog.Logger) error {

	// Setup OAuth providers
	goth.UseProviders(
		google.New(cfg.OAuth.ClientID, cfg.OAuth.ClientSecret, cfg.OAuth.ClientCallbackURL, "email", "profile"),
	)

	// Set gin to release mode so we get clean logs
	gin.SetMode(gin.ReleaseMode)

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.ClientFly, cfg.ClientLocal, cfg.ClientProxy},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Initialise gin router
	router := gin.New()
	router.Use(cors.New(corsConfig)) // pass cors config to gin router

	// This means all our logs will be same format instead of a mix between gins and slogs
	router.Use(CustomLogger(logger))

	// Register routes
	routes.Routes(router, db, logger)

	// Start the server
	logger.Info("Starting Server", "port", cfg.Port)
	err := router.Run(("0.0.0.0:" + cfg.Port))

	return err
}
