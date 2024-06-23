package routes

import (
	"fmt"
	"gymBackend/models"
	"gymBackend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSession(context *gin.Context) {

	token := context.GetHeader("Authorization")

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var session models.Session
	if err := context.ShouldBindJSON(&session); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercises []models.Exercise

	if err := session.CreateSession(exercises, userId); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Print received session data for debugging
	fmt.Printf("Received session: %+v\n", session)

	context.JSON(http.StatusOK, gin.H{"message": "Session created successfully"})
}

func CheckInSession(context *gin.Context) {
	var sessionID struct {
		ID int `json:"id"`
	}
	if err := context.ShouldBindJSON(&sessionID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := models.MarkSessionCompleted(sessionID.ID); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Status": "checked in"})
}

func GetSession(context *gin.Context) {
	sessionIDStr := context.Param("id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	session, err := models.GetSession(sessionID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, session)
}
