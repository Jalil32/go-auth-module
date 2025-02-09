package auth

import (
	"fmt"
	"net/http"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// SignInWithProvider handles third-party sign-in using a provider (e.g., Google)
func (a *AuthController) SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Provider not specified", nil)
		return
	}

	// Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Begin OAuth flow
	a.Logger.Info("Starting OAuth flow", "provider", provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler handles the OAuth callback and user creation
func (a *AuthController) CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Provider not specified", nil)
		return
	}

	// Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		a.Logger.Error("OAuth complete error", "provider", provider, "error", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OAuth complete error: %w", err))
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "OAuth complete error", err)
		return
	}

	// Check if the user exists in the database
	existingUser, err := a.UserDB.FindUserByEmail(oauthUser.Email)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Database error during lookup", err)
		return
	}

	// 4) Start transaction
	tx, err := a.UserDB.Beginx()
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to start transaction", err)
		return
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.Logger.Error("Failed to rollback transaction", "error", err)
			}
		}
	}()

	var user *models.User
	if existingUser == nil {
		newUser := models.User{
			Email:     oauthUser.Email,
			FirstName: oauthUser.FirstName,
			LastName:  oauthUser.LastName,
			Provider:  &oauthUser.Provider,
			Verified:  true,
		}
		err := a.UserDB.CreateUser(tx, &newUser)
		if err != nil {
			a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to create user", err)
			return
		}
		user = &newUser
	} else {
		user = existingUser
	}

	// Generate JWT token
	token, err := a.JWTGenerator.GenerateJWT(user)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate JWT", err)
		return
	}

	// Set the JWT token as a cookie
	a.setAuthCookie(c, token)

	// 7) Commit the transaction
	if err := tx.Commit(); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Faile to commit transaction", err)
		return
	}

	// Redirect to the /dashboard page
	c.Redirect(http.StatusFound, a.FrontendAddress+"/auth/otp")
}
