package middleware

import (
	"log/slog"
)

type Middleware struct {
	Logger *slog.Logger
}

// NewAuthController initializes a new AuthController
func NewMiddlewareSetup(logger *slog.Logger) *Middleware {

	return &Middleware{
		Logger: logger,
	}
}
