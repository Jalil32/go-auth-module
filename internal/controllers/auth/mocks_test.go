package auth_test

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/jalil32/go-auth-module/internal/models"
)

// MockRedisClient is a mock implementation of the RedisClient interface.
type MockRedisClient struct {
	GetFunc func(ctx context.Context, key string) *redis.StringCmd
	SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	DelFunc func(ctx context.Context, keys ...string) *redis.IntCmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return redis.NewStringResult("", errors.New("not implemented"))
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	return redis.NewStatusResult("", errors.New("not implemented"))
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, keys...)
	}
	return redis.NewIntResult(0, errors.New("not implemented"))
}

// MockDB is a mock implementation of the UserRepository interface.
type MockDB struct {
	FindUserByEmailFunc func(email string) (*models.User, error)
	CreateUserFunc      func(ext sqlx.Ext, user *models.User) error
	UpdateUserFunc      func(ext sqlx.Ext, user *models.User) error
	BeginxFunc          func() (*sqlx.Tx, error)
}

func (m *MockDB) FindUserByEmail(email string) (*models.User, error) {
	return m.FindUserByEmailFunc(email)
}

func (m *MockDB) CreateUser(ext sqlx.Ext, user *models.User) error {
	return m.CreateUserFunc(ext, user)
}

func (m *MockDB) UpdateUser(ext sqlx.Ext, user *models.User) error {
	return m.UpdateUserFunc(ext, user)
}

func (m *MockDB) Beginx() (*sqlx.Tx, error) {
	return nil, nil
}

// MockJWTGenerator is a mock implementation of the JWTGenerator interface.
type MockJWTGenerator struct {
	GenerateJWTFunc func(user *models.User) (string, error)
}

func (m *MockJWTGenerator) GenerateJWT(user *models.User) (string, error) {
	if m.GenerateJWTFunc != nil {
		return m.GenerateJWTFunc(user)
	}
	return "", nil
}

// MockLogger is a mock implementation of the Logger interface.
type MockLogger struct {
	ErrorFunc func(msg string, keysAndValues ...interface{})
	InfoFunc  func(msg string, keysAndValues ...interface{})
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, keysAndValues...)
	}
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, keysAndValues...)
	}
}
