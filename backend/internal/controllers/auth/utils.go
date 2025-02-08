package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/go-gomail/gomail"
	"golang.org/x/crypto/bcrypt"
)

func (a *AuthController) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.Logger.Error("Failed to hash password", "error", err)
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
	// 1) Generate otp
	otp, err := a.generateOTP()
	if err != nil {
		a.Logger.Error("Failed to generate otp", "error", err)
		return fmt.Errorf("Failed to generate otp: %w", err)
	}

	// 2) Store otp
	err = a.storeOTP(email, otp)
	if err != nil {
		a.Logger.Error("Failed to store otp", "error", err)
		return fmt.Errorf("Failed to store otp: %w", err)
	}

	// 3) Create new message
	message := gomail.NewMessage()

	// 4) Set email headers
	message.SetHeader("From", "team@demomailtrap.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "WealthScope One Time Password")
	message.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s", otp))

	// 5) Send the email
	port, err := strconv.Atoi(a.Port)
	if err != nil {
		a.Logger.Error("Failed to convert port to int", "error", err)
		return fmt.Errorf("Failed to convert port: %w", err)
	}

	dailer := gomail.NewDialer(a.Host, port, a.Username, a.Password)

	if err := dailer.DialAndSend(message); err != nil {
		a.Logger.Error("Failed to send otp", "error", err)
		return fmt.Errorf("Failed to send otp: %w", err)
	}

	return nil
}
