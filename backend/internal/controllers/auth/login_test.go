package auth_test

import (
	"log/slog"
	"testing"
	"wealthscope/backend/internal/controllers/auth"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type MockDB struct{}

func (m *MockDB) 

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db     *sqlx.DB
		rdb    *redis.Client
		logger *slog.Logger
		// Named input parameters for target function.
		c *gin.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := auth.NewAuthController(tt.db, tt.rdb, tt.logger)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			a.Login(tt.c)
		})
	}
}
