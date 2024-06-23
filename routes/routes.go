package routes

import (
	"gymBackend/middleware"
	"gymBackend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", CreateUser)
	server.POST("/login", Login)

	authGroupUsers := server.Group("/users")
	authGroupUsers.Use(middleware.AuthMiddleware())
	{
		authGroupUsers.PUT("/updateUserProfile", UpdateUserProfile)
		authGroupUsers.GET("/:id", getUser)
	}

	// Admin routes group protected by SuperUserMiddleware
	adminGroupUsers := server.Group("/users")
	adminGroupUsers.Use(middleware.SuperUserMiddleware())
	{
		adminGroupUsers.GET("/showusers", ShowUsers)

	}

	adminGroupTrainers := server.Group("/trainers")
	adminGroupTrainers.Use(middleware.SuperUserMiddleware())
	{
		adminGroupTrainers.POST("/createtrainer", CreateTrainer)
		adminGroupTrainers.GET("/showtrainers", ShowTrainers)
		adminGroupTrainers.GET("/:id", getTrainer)
		adminGroupTrainers.PUT("/updateTrainerProfile", UpdateTrainerProfile)
	}

	// Routes

	authGroupSessions := server.Group("/session")
	authGroupSessions.Use(middleware.AuthMiddleware())
	{
		authGroupSessions.POST("/sessions", CreateSession)
		authGroupSessions.POST("/check_in", CheckInSession)
	}

	// Start cleanup job for missed sessions
	go func() {
		for {
			if err := models.MarkMissedSessions(); err != nil {
				// Handle error as needed
			}
			time.Sleep(5 * time.Minute) // Run every 5 minutes
		}
	}()

}
