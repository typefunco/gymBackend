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
	adminGroupUsers := server.Group("/users")
	adminGroupUsers.Use(middleware.SuperUserMiddleware())
	{
		adminGroupUsers.GET("/showusers", ShowUsers)
		adminGroupUsers.POST("/createtrainer", CreateTrainer)

	}

	adminGroupTrainers := server.Group("/trainers")
	adminGroupTrainers.Use(middleware.SuperUserMiddleware())
	{
		adminGroupTrainers.POST("/createtrainer", CreateTrainer)
		adminGroupTrainers.GET("/showtrainers", ShowTrainers)
		adminGroupTrainers.POST("/updateTrainerProfile", UpdateTrainerProfile)
	}

}
