# Go Authentication Module

A production-ready authentication microservice built with Go, featuring JWT-based session management, email verification, OAuth integration, and comprehensive password management.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Environment Variables](#environment-variables)
- [Database Migrations](#database-migrations)
- [Testing](#testing)
- [Security Features](#security-features)
- [Project Structure](#project-structure)

## Features

### Core Authentication
- **User Registration** - Email/password signup with strong password validation
- **User Login** - Secure authentication with bcrypt password hashing
- **JWT Sessions** - Stateless authentication with HTTP-only secure cookies
- **Logout** - Session termination with cookie cleanup

### Email Verification
- **OTP System** - 6-digit one-time passwords with 5-minute expiration
- **SMTP Integration** - Automated email delivery for verification codes
- **Redis-backed Storage** - Fast, ephemeral storage for OTPs

### Password Management
- **Password Reset Flow** - Secure email-based password reset
- **Reset Tokens** - UUID-based one-time use tokens with 15-minute expiration
- **Strong Password Validation** - Enforced complexity requirements:
  - Minimum 8 characters
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one digit
  - At least one special character (@$!%*?&)

### OAuth Integration
- **Google OAuth** - Social sign-in with automatic account creation
- **Provider Validation** - Extensible architecture for additional providers
- **Auto-verification** - OAuth users are automatically verified

### Security
- **Bcrypt Hashing** - Industry-standard password encryption
- **HTTP-only Cookies** - Protection against XSS attacks
- **Secure Flags** - HTTPS enforcement in production
- **Input Validation** - Comprehensive request validation with detailed error messages
- **Database Transactions** - ACID compliance for data integrity
- **Token Expiration** - Configurable JWT and reset token lifetimes

## Tech Stack

- **Language:** Go 1.22+
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **Cache/Session Store:** Redis
- **Authentication:** JWT (HS256)
- **OAuth:** Goth (Google provider)
- **Email:** SMTP via Mailtrap
- **Password Hashing:** bcrypt
- **Validation:** go-playground/validator
- **Testing:** Testify
- **Logging:** slog with tint formatter
- **Migrations:** Goose

## Architecture

### Design Patterns
- **Repository Pattern** - Clean separation of data access logic
- **Dependency Injection** - Testable, loosely-coupled components
- **Interface-based Design** - Mock-friendly architecture for testing
- **Transaction Management** - Atomic operations with rollback support

### Security Architecture
- **Defense in Depth** - Multiple layers of security controls
- **Stateless Sessions** - JWT-based authentication for horizontal scaling
- **Token Rotation** - One-time use tokens for sensitive operations
- **Secure Communication** - HTTPS enforcement, secure cookies

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 13+
- Redis 6+
- Docker & Docker Compose (for development database)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/jalil32/go-auth-module.git
   cd go-auth-module
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.template .env
   ```

   Edit `.env` and configure the following:
   - Database credentials (PostgreSQL)
   - Redis connection details
   - JWT secret and expiration
   - SMTP settings (Mailtrap)
   - Google OAuth credentials (optional)

4. **Generate JWT secret**
   ```bash
   # Linux/macOS
   openssl rand -base64 32

   # Add the output to your .env as JWT_TOKEN
   ```

5. **Start development database**
   ```bash
   cd deployments
   docker-compose up -d
   ```

6. **Start Redis**
   ```bash
   # Follow installation for your OS: https://redis.io/docs/latest/operate/oss_and_stack/install/
   redis-server
   ```

7. **Run database migrations**
   ```bash
   export $(cat .env | xargs)
   goose -dir ./migrations postgres "user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_NAME host=$POSTGRES_HOST port=$POSTGRES_PORT sslmode=$POSTGRES_SSL_MODE" up
   ```

8. **Run the application**
   ```bash
   go run cmd/app/main.go
   ```

The server will start on `http://localhost:3000`

## API Documentation

### Base URL
```
http://localhost:3000/api
```

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response** (201 Created):
```json
{
  "message": "User created and OTP sent successfully"
}
```

---

#### Verify OTP
```http
POST /api/auth/verify
Content-Type: application/json

{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Response** (201 Created):
```json
{
  "message": "OTP verified successfully"
}
```
*Sets `auth_token` cookie*

---

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response** (200 OK):
```json
{
  "message": "Login successful"
}
```
*Sets `auth_token` cookie*

---

#### Logout
```http
POST /api/auth/logout
```

**Response** (200 OK):
```json
{
  "message": "Logout successful"
}
```
*Clears `auth_token` cookie*

---

#### Forgot Password
```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response** (200 OK):
```json
{
  "message": "If a user with this email exists, a password reset link has been sent."
}
```
*Sends reset email with token*

---

#### Reset Password
```http
POST /api/auth/reset-password?token=<reset-token>
Content-Type: application/json

{
  "newPassword": "NewSecurePass123!"
}
```

**Response** (200 OK):
```json
{
  "message": "Password reset successfully."
}
```

---

#### OAuth Sign-In (Google)
```http
GET /api/auth/google
```

**Response**: Redirects to Google OAuth consent screen

---

#### OAuth Callback
```http
GET /api/auth/google/callback
```

**Response**: Redirects to `/dashboard` with `auth_token` cookie set

---

### Protected Routes

All protected routes require the `auth_token` cookie or `Authorization` header.

#### Example Protected Endpoint
```http
GET /protected
Cookie: auth_token=<jwt-token>
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

**Response** (401 Unauthorized):
```json
{
  "error": "Missing token"
}
```

---

## Environment Variables

Create a `.env` file in the config directory with the following variables:

```bash
# PostgreSQL Configuration
POSTGRES_USER=admin
POSTGRES_PASSWORD=your_secure_password
POSTGRES_NAME=auth_db
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_SSL_MODE=disable

# Redis Configuration
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# JWT Configuration
JWT_TOKEN=your_generated_secret_key_here
JWT_EXPIRY=1h

# SMTP Configuration (Mailtrap)
EMAIL_HOST=live.smtp.mailtrap.io
EMAIL_PORT=587
EMAIL_USERNAME=smtp@mailtrap.io
EMAIL_PASSWORD=your_mailtrap_password

# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CLIENT_CALLBACK_URL=http://localhost:3000/api/auth/google/callback

# Server Configuration
BACKEND_PORT=3000
CLIENT_LOCAL=http://localhost:5173
BACKEND_LOCAL=http://localhost:3000
CLIENT_PROXY=http://localhost:3000

# Production (Fly.io)
CLIENT_FLY=https://your-client-app.fly.dev
```

## Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

### Creating a Migration
```bash
cd migrations
goose create add_users_table sql
```

### Running Migrations
```bash
export $(cat .env | xargs)
goose -dir ./migrations postgres "user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_NAME host=$POSTGRES_HOST port=$POSTGRES_PORT sslmode=$POSTGRES_SSL_MODE" up
```

### Rolling Back
```bash
export $(cat .env | xargs)
goose -dir ./migrations postgres "user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_NAME host=$POSTGRES_HOST port=$POSTGRES_PORT sslmode=$POSTGRES_SSL_MODE" down
```

## Testing

### Run Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Specific Test
```bash
go test -run TestAuthController_Login ./internal/controllers/auth
```

### Test Coverage Highlights
- Login flow with multiple test cases
- Mock implementations for database, Redis, JWT, and logger
- Edge case handling (missing credentials, invalid users, database errors)

## Security Features

### Password Security
- **Bcrypt Hashing**: All passwords hashed with bcrypt.DefaultCost (10 rounds)
- **Strong Password Requirements**: Enforced complexity rules
- **No Password Logging**: Passwords never appear in logs

### Session Security
- **HTTP-only Cookies**: JavaScript cannot access auth tokens
- **Secure Flag**: HTTPS-only in production
- **Configurable Expiration**: Token lifetime controlled via environment variable

### Token Security
- **HMAC Signature**: HS256 signing algorithm
- **Expiration Validation**: Automatic token expiry checking
- **One-time Use**: Password reset tokens deleted after use

### Data Protection
- **Database Transactions**: Atomic operations with rollback on failure
- **Input Validation**: All requests validated before processing
- **SQL Injection Prevention**: Parameterized queries via sqlx

### Anti-Enumeration
- **Consistent Responses**: Same message for existing/non-existing users in password reset
- **Generic Error Messages**: User-friendly errors without sensitive details

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go                 # Application entry point
├── config/
│   └── env.go                      # Configuration loader
├── internal/
│   ├── app/
│   │   ├── server.go               # Server initialization
│   │   └── gin_custom_logger.go   # Custom Gin logger
│   ├── controllers/
│   │   └── auth/
│   │       ├── auth_controller.go  # Controller initialization
│   │       ├── register.go         # Registration handler
│   │       ├── login.go            # Login handler
│   │       ├── logout.go           # Logout handler
│   │       ├── oauth.go            # OAuth handlers
│   │       ├── otp.go              # OTP verification handler
│   │       ├── otp_util.go         # OTP generation & sending
│   │       ├── forgot_password.go  # Password reset handlers
│   │       ├── jwt_util.go         # JWT generation
│   │       ├── password_util.go    # Password hashing
│   │       ├── validator_util.go   # Request validation
│   │       └── error.go            # Error handling
│   ├── db/
│   │   ├── connection.go           # Database connection
│   │   └── user_repository.go     # User data access
│   ├── middleware/
│   │   ├── auth_middleware.go      # JWT validation middleware
│   │   └── logger_middleware.go    # Request logging
│   ├── models/
│   │   └── user_model.go          # User data model
│   └── routes/
│       └── routes.go              # Route definitions
├── migrations/                     # Database migrations
├── deployments/
│   └── compose.yml                # Docker Compose for dev DB
├── .env.template                  # Environment variables template
├── fly.toml                       # Fly.io deployment config
├── go.mod                         # Go module dependencies
└── README.md                      # This file
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Contact

Jalil - [@jalil32](https://github.com/jalil32)

Project Link: [https://github.com/jalil32/go-auth-module](https://github.com/jalil32/go-auth-module)

---

**Note:** This is a portfolio/demonstration project showcasing production-ready authentication patterns in Go. While the core authentication features are fully implemented and tested, some additional features (like rate limiting and account lockout) are not yet implemented but would be recommended for a complete production deployment.
