package middlewares

import (
	"net/http"

	"github.com/asdutoit/gotraining/section11/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")

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
