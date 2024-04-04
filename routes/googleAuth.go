package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asdutoit/go_backend_template/models"
	"github.com/asdutoit/go_backend_template/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func handleGoogleAuth(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No id_token field in oauth2 token."})
		return
	}

	// You can now use the rawIDToken to authenticate the user in your system.
	idToken, err := idtoken.Validate(context.Background(), rawIDToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate ID token"})
		return
	}

	// The ID token contains the user's profile information in the Claims field.
	email := idToken.Claims["email"]
	givenName := idToken.Claims["given_name"]
	familyName := idToken.Claims["family_name"]
	picture := idToken.Claims["picture"]

	// 1.  Check if user exists in database
	user, err := models.GetUserByEmail(email.(string))

	// Check for database errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user", "error": err.Error()})
		return
	}

	if user == nil {
		user = &models.User{
			Email:      email.(string),
			First_name: givenName.(string),
			Last_name:  familyName.(string),
			Picture:    picture.(string),
			Username:   email.(string),
		}
		err = user.Save()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user", "error": err.Error()})
			return
		}

		jwtToken, err := utils.GenerateToken(user.Email, user.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token", "error": err.Error()})
			return
		}

		RedirectWithCookie(c, jwtToken, token.Expiry, user)
	} else {
		user.First_name = givenName.(string)
		user.Last_name = familyName.(string)
		user.Picture = picture.(string)

		err = user.Update()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": err.Error()})
			return
		}

		jwtToken, err := utils.GenerateToken(user.Email, user.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token", "error": err.Error()})
			return
		}
		RedirectWithCookie(c, jwtToken, token.Expiry, user)
	}
}

func RedirectWithCookie(ctx *gin.Context, jwtToken string, expiry time.Time, user *models.User) {
	// Determine whether we're in a secure environment
	secureEnv := os.Getenv("SECURE_ENV") == "true"

	// Determine the SameSite mode based on the environment
	var sameSiteMode http.SameSite
	if secureEnv {
		sameSiteMode = http.SameSiteStrictMode
	} else {
		sameSiteMode = http.SameSiteLaxMode
	}

	claims := jwt.MapClaims{
		"userId":     user.ID,
		"first_name": user.First_name,
		"last_name":  user.Last_name,
		"email":      user.Email,
		"picture":    user.Picture,
		"exp":        expiry.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Printf("Err creating token: %v", err)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    ss,
		Expires:  expiry,
		HttpOnly: true,
		Secure:   secureEnv,
		SameSite: sameSiteMode,
		Path:     "/",
	})

	ctx.Redirect(http.StatusMovedPermanently, os.Getenv("FRONTEND_URL"))
}
