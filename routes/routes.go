package routes

import (
	"github.com/asdutoit/go_backend_template/graphql"
	"github.com/asdutoit/go_backend_template/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))

	// Add your routes here
	r.GET("healthcheck", healthCheck)
	r.POST("/signup", signUp)
	r.POST("/login", login)
	r.POST("/events", middlewares.Authenticate, createEvent)
	r.GET("/deployments", getDeployments)
	r.GET("/deployment", GetDeploymentByQuery)
	r.GET("/auth/google", handleGoogleAuth)
	r.GET("/auth/google/callback", handleGoogleCallback)

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
	r.Any("/graphql", func(c *gin.Context) {
		h := graphql.NewHandler()
		h.ContextHandler(c.Request.Context(), c.Writer, c.Request)
	})

	return r
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
