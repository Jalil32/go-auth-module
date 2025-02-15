package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *AuthController) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (a *AuthController) storeOTP(email string, otp string) error {
	ctx := context.Background()
	return a.RedisCache.Set(ctx, email, otp, 5*time.Minute).Err() // Store OTP for 5 minutes
}

func (a *AuthController) validateOTP(email string, inputOTP string) bool {
	ctx := context.Background()
	storedOTP, err := a.RedisCache.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	return storedOTP == inputOTP
}

func (a *AuthController) generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // Range: 0 - 999999
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil // Ensure it’s always 6 digits
}

func (a *AuthController) sendOTP(email string) error {
	// 1) Generate OTP
	otp, err := a.generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate otp: %w", err)
	}

	// 2) Store OTP
	err = a.storeOTP(email, otp)
	if err != nil {
		return fmt.Errorf("failed to store otp: %w", err)
	}

	// 3) Create new message
	message := gomail.NewMessage()

	// 4) Set email headers
	message.SetHeader("From", "team@demomailtrap.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "WealthScope One Time Password")
	message.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s", otp))

	// 5) Convert port to int
	port, err := strconv.Atoi(a.Port)
	if err != nil {
		return fmt.Errorf("failed to convert port: %w", err)
	}

	// 6) Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 7) Create a channel to handle the result of the email sending
	done := make(chan error, 1)

	// 8) Run the email sending in a goroutine
	go func() {
		dialer := gomail.NewDialer(a.Host, port, a.Username, a.Password)
		done <- dialer.DialAndSend(message)
	}()

	// 9) Wait for either the email to be sent or the context to timeout
	select {
	case <-ctx.Done():
		// Context timed out
		return fmt.Errorf("failed to send OTP: timeout reached")
	case err := <-done:
		// Email sending completed
		if err != nil {
			return fmt.Errorf("failed to send OTP: %w", err)
		}
	}

	return nil
}

func (a *AuthController) sendForgotPasswordToken(email string, link string) error {
	// 1) Create new message
	message := gomail.NewMessage()

	// 2) Set email headers
	message.SetHeader("From", "team@demomailtrap.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Forgot Password")
	message.SetBody("text/plain", fmt.Sprintf("Please click the link to reset your password: %s", link))

	// 3) Convert port to int
	port, err := strconv.Atoi(a.Port)
	if err != nil {
		return fmt.Errorf("failed to convert port: %w", err)
	}

	// 4) Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5) Create a channel to handle the result of the email sending
	done := make(chan error, 1)

	// 6) Run the email sending in a goroutine
	go func() {
		dialer := gomail.NewDialer(a.Host, port, a.Username, a.Password)
		done <- dialer.DialAndSend(message)
	}()

	// 7) Wait for either the email to be sent or the context to timeout
	select {
	case <-ctx.Done():
		// Context timed out
		return fmt.Errorf("failed to send forgot password link: timeout reached")
	case err := <-done:
		// Email sending completed
		if err != nil {
			return fmt.Errorf("failed to send forgot password link: %w", err)
		}
	}

	return nil
}

func (a *AuthController) generateForgotPasswordLink(email string) (string, error) {
	token := uuid.New().String()

	key := fmt.Sprintf("forgot_password:%v", token)
	ctx := context.Background()
	err := a.RedisCache.Set(ctx, key, email, 15*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("Storing token in redis failed: %v", err)
	}

	link := fmt.Sprintf("%s/reset-password?token=%s", a.FrontendAddress, token)
	return link, nil
}

// Validate forgot password token
func (a *AuthController) validateForgotPasswordToken(token string) (string, error) {
	ctx := context.Background()

	key := fmt.Sprintf("forgot_password:%v", token)
	email, err := a.RedisCache.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("Invalid or expired token: %v", err)
	}

	a.RedisCache.Del(ctx, key)
	return email, nil
}
