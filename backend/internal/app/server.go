package server

import (
	"log/slog"
	"wealthscope/backend/config"
	"wealthscope/backend/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func StartServer(cfg *config.Config, db *sqlx.DB, logger *slog.Logger) error {

	// Initalise gin router
	router := gin.Default()

	// Add authentication routes
	routes.Routes(router)

	// Start the server
	logger.Info("Starting Server", "port", cfg.Port)
	err := router.Run((":" + cfg.Port))

	return err
}
