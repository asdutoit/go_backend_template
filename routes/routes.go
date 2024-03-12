package routes

import (
	"github.com/asdutoit/gotraining/section11/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// Add your routes here
	r.GET("healthcheck", healthCheck)
	r.POST("/signup", signUp)
	r.POST("/login", login)
	r.POST("/events", middlewares.Authenticate, createEvent)

	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.GET("/events", getEvents)
	authenticated.GET("/events/:id", getEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.GET("/user", getUsers)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)
	authenticated.POST("/uploads", uploadFiles)
	authenticated.POST("/deleteUser", deleteUser)

	return r
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
