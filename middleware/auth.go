package middleware

import (
	"fmt"
	"gymBackend/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func jwtSecret() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	secretKey := os.Getenv("secretKey")
	return secretKey
}

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")

	if tokenString == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"response": "Authorization header is required"})
		return
	}

	userId, err := utils.VerifyToken(tokenString)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"response": "Wrong authorization header"})
		return
	}

	context.Set("userId", userId)
	context.Next()
}
