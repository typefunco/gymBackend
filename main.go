package main

import (
	// "gymBackend/models"
	"gymBackend/routes"
	"gymBackend/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/withmandala/go-log"
)

func main() {
	logger := log.New(os.Stderr)
	connection := utils.CheckConnection()
	if !connection {
		logger.Fatal("CAN'T CONNECT TO DB")
		return
	}
	logger.Info("Connected to DB")

	server := gin.Default()

	routes.RegisterRoutes(server) // Setted up server routes
	server.Run()
}
