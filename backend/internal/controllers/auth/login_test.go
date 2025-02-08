package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"wealthscope/backend/config"
	"wealthscope/backend/internal/controllers/auth"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// MockRedisClient is a mock implementation of the RedisClient interface.
type MockRedisClient struct {
	GetFunc func(ctx context.Context, key string) *redis.StringCmd
	SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
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

// MockDB is a mock implementation of the UserRepository interface.
type MockDB struct {
	FindUserByEmailFunc func(email string) (*models.User, error)
	CreateUserFunc      func(ext sqlx.Ext, user *models.User) error
	UpdateUserFunc      func(ext sqlx.Ext, user *models.User) error
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

// Helper function to create a test AuthController.
func createTestAuthController(
	mockDB *MockDB,
	mockRedis *MockRedisClient,
	mockLogger *MockLogger,
	mockJWTGenerator *MockJWTGenerator,
) (*auth.AuthController, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return auth.NewAuthController(mockDB, mockRedis, mockLogger, mockJWTGenerator, cfg)
}

// Helper function to create a test HTTP request.
func createTestRequest(method, url string, body map[string]string) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Helper function to execute the Login handler and return the response.
func executeLoginHandler(authController *auth.AuthController, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	authController.Login(c)
	return w
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_EXPIRY", "1m")

	// Generate a static bcrypt password hash for consistency in tests.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPasswordStr := string(hashedPassword)
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockDB         *MockDB
		mockJWT        *MockJWTGenerator
		mockLogger     *MockLogger
		mockRedis      *MockRedisClient
		expectedStatus int
		expectedBody   string
	}{
		// {
		// 	name:        "Empty Request Body",
		// 	requestBody: map[string]string{},
		// 	mockDB:      &MockDB{},
		// 	mockJWT:     &MockJWTGenerator{},
		// 	mockLogger: &MockLogger{
		// 		ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
		// 	},
		// 	mockRedis:      &MockRedisClient{},
		// 	expectedStatus: http.StatusBadRequest,
		// 	expectedBody:   `{"error":"Invalid request payload"}`,
		// },
		{
			name: "Missing Email",
			requestBody: map[string]string{
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, errors.New("database error")
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "", errors.New("JWT generation failed")
				},
			},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"details":"Key: 'requestPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag", "error":"Invalid input"}`,
		},
		{
			name: "Missing Password",
			requestBody: map[string]string{
				"email": "test@example.com",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, errors.New("database error")
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "", errors.New("JWT generation failed")
				},
			},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"details":"Key: 'requestPayload.Password' Error:Field validation for 'Password' failed on the 'required' tag", "error":"Invalid input"}`,
		},
		{
			name: "User Not Found",
			requestBody: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, nil // User not found
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "", errors.New("JWT generation failed")
				},
			},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid email or password"}`,
		},
		{
			name: "Database Error",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, errors.New("database error")
				},
			},
			mockJWT: &MockJWTGenerator{},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Database error during user lookup"}`,
		},
		// Need to mock otp service
		// {
		// 	name: "User Not Verified",
		// 	requestBody: map[string]string{
		// 		"email":    "test@example.com",
		// 		"password": "password123",
		// 	},
		// 	mockDB: &MockDB{
		// 		FindUserByEmailFunc: func(email string) (*models.User, error) {
		// 			return &models.User{
		// 				Email:        email,
		// 				PasswordHash: &hashedPasswordStr,
		// 				Verified:     false,
		// 			}, nil
		// 		},
		// 	},
		// 	mockJWT: &MockJWTGenerator{},
		// 	mockLogger: &MockLogger{
		// 		InfoFunc: func(msg string, keysAndValues ...interface{}) {},
		// 	},
		// 	mockRedis: &MockRedisClient{
		// 		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
		// 			return redis.NewStatusResult("OK", nil)
		// 		},
		// 	},
		// 	expectedStatus: http.StatusOK,
		// 	expectedBody:   `{"message":"User is not verified.","user":{"email":"test@example.com","firstName":"","lastName":""}}`,
		// },
		{
			name: "JWT Generation Failure",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return &models.User{
						Email:        email,
						PasswordHash: &hashedPasswordStr,
						Verified:     true,
					}, nil
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "", errors.New("JWT generation failed")
				},
			},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to generate token"}`,
		},
		// {
		// 	name: "Redis Set Failure",
		// 	requestBody: map[string]string{
		// 		"email":    "test@example.com",
		// 		"password": "password123",
		// 	},
		// 	mockDB: &MockDB{
		// 		FindUserByEmailFunc: func(email string) (*models.User, error) {
		// 			return &models.User{
		// 				Email:        email,
		// 				PasswordHash: &hashedPasswordStr,
		// 				Verified:     false,
		// 			}, nil
		// 		},
		// 	},
		// 	mockJWT: &MockJWTGenerator{
		// 		GenerateJWTFunc: func(user *models.User) (string, error) {
		// 			return "mock-token", nil
		// 		},
		// 	},
		// 	mockLogger: &MockLogger{
		// 		ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
		// 	},
		// 	mockRedis: &MockRedisClient{
		// 		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
		// 			return redis.NewStatusResult("", errors.New("Redis set failed"))
		// 		},
		// 	},
		// 	expectedStatus: http.StatusInternalServerError,
		// 	expectedBody:   `{"error":"Failed to send otp"}`,
		// },
		{
			name: "User Needs to Sign In with Provider",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					provider := "google"
					return &models.User{
						Email:        email,
						PasswordHash: nil,
						Provider:     &provider,
					}, nil
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "mock-token", nil
				},
			},
			mockLogger: &MockLogger{
				ErrorFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Please sign in with a provider"}`,
		},
		{
			name: "Successful Login",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return &models.User{
						Email:        email,
						PasswordHash: &hashedPasswordStr,
						Verified:     true,
					}, nil
				},
			},
			mockJWT: &MockJWTGenerator{
				GenerateJWTFunc: func(user *models.User) (string, error) {
					return "mock-token", nil
				},
			},
			mockLogger: &MockLogger{
				InfoFunc: func(msg string, keysAndValues ...interface{}) {},
			},
			mockRedis: &MockRedisClient{
				GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult("", nil)
				},
				SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
					return redis.NewStatusResult("OK", nil)
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Login successful","token":"mock-token"}`,
		},
		{
			name: "Invalid Email",
			requestBody: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, nil
				},
			},
			mockJWT:        &MockJWTGenerator{},
			mockLogger:     &MockLogger{ErrorFunc: func(msg string, keysAndValues ...interface{}) {}},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid email or password"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the AuthController.
			authController, err := createTestAuthController(tt.mockDB, tt.mockRedis, tt.mockLogger, tt.mockJWT)
			if err != nil {
				t.Fatalf("failed to create AuthController: %v", err)
			}

			// Create the HTTP request.
			req, err := createTestRequest(http.MethodPost, "/login", tt.requestBody)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Execute the Login handler.
			w := executeLoginHandler(authController, req)

			// Assertions.
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
