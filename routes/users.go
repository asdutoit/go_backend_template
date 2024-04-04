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

	usersResponse := make([]models.UserResponse, len(users))

	for i, user := range users {
		usersResponse[i] = models.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			First_name: user.First_name,
			Last_name:  user.Last_name,
			Picture:    user.Picture,
		}
	}

	ctx.JSON(http.StatusOK, usersResponse)
}

func getUser(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user ID"})
		return
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not an integer"})
		return
	}

	user, err := models.GetUserById(userIdInt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch user", "error": err.Error()})
		return
	}

	userResponse := models.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		First_name: user.First_name,
		Last_name:  user.Last_name,
		Picture:    user.Picture,
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// func getUserFromCookie(ctx *gin.Context) {
// 	// Decode user details from the token in cookie
// 	tokenCookie, err := ctx.Request.Cookie("token")
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token cookie"})
// 		return
// 	}
// 	tokenString := tokenCookie.Value

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET")), nil
// 	})

// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		email := claims["email"].(string)
// 		user, err := models.GetUserByEmail(email)

// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch user", "error": err.Error()})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, user)
// 	} else {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 	}
// }

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
