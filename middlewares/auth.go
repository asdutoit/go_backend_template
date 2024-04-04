package middlewares

import (
	"fmt"
	"net/http"

	"github.com/asdutoit/go_backend_template/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	fmt.Println(token)

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
		return
	}

	userId, err := utils.ValidateToken(token)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	ctx.Set("userId", userId)
	ctx.Next()
}
