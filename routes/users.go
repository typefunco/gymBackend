package routes

import (
	"gymBackend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/withmandala/go-log"
)

func Get(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"Message": "All done"})
}

func CreateUser(context *gin.Context) {
	logger := log.New(os.Stderr)
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": "Wrong URL address"})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": "Some problems"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"Response": "User saved"})
	logger.Infof("User saved\nFirst name: %s\nLast name: %s\nWeight: %.2f\nHeight: %.2f\nFat %%: %.2f\nMuscle %%, %.2f", user.FirstName, user.LastName, user.Weight, user.Height, user.FatPercentage, user.MusclePercentage)
}
