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

func ShowTrainers(context *gin.Context) {
	trainers, err := models.GetTrainers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Response": trainers})
}

func UpdateTrainerProfile(context *gin.Context) {

	var updates map[string]interface{}
	if err := context.BindJSON(&updates); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	trainer := &models.Trainer{}
	if err := trainer.UpdateProfile(trainer.Id, updates); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Trainer profile updated successfully"})
}
