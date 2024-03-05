package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "superSecret"

func GenerateToken(email string, userId int64) (string, error) {
	// Generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	// Return the token
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (int64, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return 0, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("invalid token")
	}

	// Return the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))
	fmt.Println("userId", userId)
	tokenIsValid := token.Valid

	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	return userId, err
}
