package routes

import (
	"gymBackend/models"
	"gymBackend/utils"
	"net/http"
	"os"

	"github.com/withmandala/go-log"

	"github.com/gin-gonic/gin"
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
		context.JSON(http.StatusBadRequest, gin.H{"Response": "CAN'T SAVE USER. USER ALREADY EXIST"})
		logger.Warn("CAN'T SAVE USER. USER ALREADY EXIST")
		return
	}

	context.JSON(http.StatusCreated, gin.H{"Response": "User saved"})
	hashedPassword, _ := utils.HashPassword(user.Password)
	logger.Infof("User saved\nUsername: %s\nPassword: %s", user.Username, hashedPassword)
}

func Login(context *gin.Context) {
	logger := log.New(os.Stderr)
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Can't parse data from request"})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Username, user.Id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
	logger.Infof("User {%s} logged in", user.Username)
}

func ShowUsers(context *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Response": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Response": users})
}

// func AddData(){
// 	logger := log.New(os.Stderr)
// 	logger.Infof("User saved\nFirst name: %s\nLast name: %s\nWeight: %.2f\nHeight: %.2f\nFat %%: %.2f\nMuscle %%, %.2f", user.FirstName, user.LastName, user.Weight, user.Height, user.FatPercentage, user.MusclePercentage)
// }
