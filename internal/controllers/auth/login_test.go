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
	"github.com/jalil32/go-auth-module/config"
	"github.com/jalil32/go-auth-module/internal/controllers/auth"
	"github.com/jalil32/go-auth-module/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

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
	// 1) Set gin to test mode and set environment variables
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_EXPIRY", "1m")

	// 2) Generate a static bcrypt password hash for consistency in tests.
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
			expectedBody:   `{"error":"Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag", "internal_message":"Validation failed", "message":"Please correct the following errors:\n- The Email field is required.\n"}`,
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
			expectedBody:   `{"error":"Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", "internal_message":"Validation failed", "message":"Please correct the following errors:\n- The Password field is required.\n"}`,
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
			expectedBody:   `{"error":"User not found", "internal_message":"User not found", "message":"Invalid email or password"}`,
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
			expectedBody:   `{"error":"database error", "internal_message":"Database lookup error", "message":"Something went wrong..."}`,
		},
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
			expectedBody:   `{"error":"JWT generation failed", "internal_message":"Failed to generate JWT Token", "message":"Something went wrong..."}`,
		},
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
			expectedBody:   `{"error":"User needs to sign in with provider", "internal_message":"User needs to sign in with provider", "message":"Please sign in with a provider"}`,
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
			expectedBody:   `{"message":"Login successful"}`,
		},
		{
			name: "Invalid Email",
			requestBody: map[string]string{
				"email":    "nonexistentexample.com",
				"password": "Password123!",
			},
			mockDB: &MockDB{
				FindUserByEmailFunc: func(email string) (*models.User, error) {
					return nil, nil
				},
			},
			mockJWT:        &MockJWTGenerator{},
			mockLogger:     &MockLogger{ErrorFunc: func(msg string, keysAndValues ...interface{}) {}},
			mockRedis:      &MockRedisClient{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag", "internal_message":"Validation failed", "message":"Please correct the following errors:\n- The Email field must be a valid email address.\n"}`,
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
