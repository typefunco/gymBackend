package main

import (
	"gymBackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server) // Setted up server routes
	server.Run()
}
