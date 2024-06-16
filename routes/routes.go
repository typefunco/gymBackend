package routes

import (
	"gymBackend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", CreateUser)
	server.POST("/login", Login)
	server.POST("/updateUserProfile", UpdateUserProfile)

	// Admin routes group protected by SuperUserMiddleware
	adminGroup := server.Group("/users")
	adminGroup.Use(middleware.SuperUserMiddleware())
	{
		adminGroup.GET("/showusers", ShowUsers)

	}
}
