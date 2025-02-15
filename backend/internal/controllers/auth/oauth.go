package auth

import (
	"net/http"
	"wealthscope/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// SignInWithProvider handles third-party sign-in using a provider (e.g., Google)
func (a *AuthController) SignInWithProvider(c *gin.Context) {
	// 1) Check that provider exists
	provider := c.Param("provider")
	if provider == "" {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Provider not specified", nil)
		return
	}

	// 2) Validate provider
	allowedProviders := map[string]bool{
		"google": true,
	}

	if !allowedProviders[provider] {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Invalid provider specified", nil)
		return
	}

	// 3) Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// 4) Begin OAuth flow
	a.Logger.Info("Starting OAuth flow", "provider", provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *AuthController) CallbackHandler(c *gin.Context) {
	// 1) Check that the provider exists
	provider := c.Param("provider")
	if provider == "" {
		a.HandleError(c, http.StatusBadRequest, "Bad Request", "Provider not specified", nil)
		return
	}

	// 2) Add provider to the request URL
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "OAuth complete error", err)
		return
	}

	// 3) Check if the user exists in the database
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

	// 5) If user does not exist, create user and set verified to true
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

	// 6) Generate JWT token
	token, err := a.JWTGenerator.GenerateJWT(user)
	if err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Failed to generate JWT", err)
		return
	}

	// 7) Set the JWT token as a cookie
	a.setAuthCookie(c, token)

	// 8) Commit the transaction
	if err := tx.Commit(); err != nil {
		a.HandleError(c, http.StatusInternalServerError, "Something went wrong...", "Faile to commit transaction", err)
		return
	}

	// 9) Redirect to the /dashboard page
	c.Redirect(http.StatusFound, a.FrontendAddress+"/dashboard")
}
