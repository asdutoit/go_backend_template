package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asdutoit/gotraining/section11/models"
	"github.com/asdutoit/gotraining/section11/utils"
	"github.com/gin-gonic/gin"
)

func getEvents(ctx *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

func createEvent(ctx *gin.Context) {
	var event models.Event
	err := ctx.ShouldBindJSON(&event)

	userId := ctx.GetInt64("userId")
	event.UserID = userId

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = event.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not create event", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "created", "event": event})
}

func getEvent(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := models.GetEventByID(eventId)
	userId := ctx.GetInt64("userId")
	event.UserID = userId
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, event)
}

func updateEvent(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userId := ctx.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	fmt.Println("userId", userId)
	fmt.Println("event user", event.UserID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "event not found"})
		return
	}

	if event.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var updatedEvent models.Event
	err = ctx.ShouldBindJSON(&updatedEvent)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.UserID = userId
	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not update event", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "Event Updated Successfully", "event": updatedEvent})
}

func deleteEvent(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
		return
	}

	_, err := utils.ValidateToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "event not found"})
		return
	}

	err = event.Delete()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete event", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "Event Deleted Successfully"})
}
