package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/get", Get)
	server.POST("/createuser", CreateUser)
	server.POST("/login", Login)
	server.GET("/showusers", ShowUsers)
}
