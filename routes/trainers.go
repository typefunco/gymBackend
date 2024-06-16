package routes

import (
	"gymBackend/models"
	"net/http"
	"os"

	"github.com/withmandala/go-log"

	"github.com/gin-gonic/gin"
)

func CreateTrainer(context *gin.Context) {
	logger := log.New(os.Stderr)
	var trainer models.Trainer
	err := context.ShouldBindJSON(&trainer)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": "Wrong URL address"})
		return
	}

	err = trainer.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": "CAN'T SAVE TRAINER"})
		logger.Warn("CAN'T SAVE TRAINER")
		return
	}

	context.JSON(http.StatusCreated, gin.H{"Response": "Trainer saved"})
	logger.Infof("Trainer created\nName: %s\nLast name: %s", trainer.FirstName, trainer.LastName)
}