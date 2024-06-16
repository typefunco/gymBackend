package middleware

import (
	"gymBackend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")

		if tokenString == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authorization header is required"})
			return
		}

		userId, err := utils.VerifyToken(tokenString)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Wrong authorization header"})
			return
		}

		context.Set("userId", userId)
		context.Next()
	}
}
