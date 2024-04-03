package routes

import (
	"net/http"

	"github.com/asdutoit/go_backend_template/models"
	"github.com/asdutoit/go_backend_template/utils"
	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {
	// Create a new user
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Save()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func getUsers(ctx *gin.Context) {
	users, err := models.GetAllUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func login(ctx *gin.Context) {
	// Declare the user variable
	var user models.User
	// Bind the request body to the user variable
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := user.ValidateCredentials()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials", "error": err.Error()})
		return
	}
	// Call the login method on the user
	token, err := utils.GenerateToken(user.Email, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func deleteUser(ctx *gin.Context) {
	// Declare the user variable
	var user models.User
	// Bind the request body to the user variable
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = user.ValidateCredentials()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials", "error": err.Error()})
		return
	}

	err = user.DeleteUser()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "User Deleted Successfully"})
}
