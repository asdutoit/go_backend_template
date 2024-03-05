package routes

import (
	"net/http"
	"strconv"

	"github.com/asdutoit/gotraining/section11/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event", "error": err.Error()})
		return
	}

	err = event.Register(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not register for event", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func cancelRegistration(ctx *gin.Context) {
	userID := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel registration", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "Cancelled"})

}
