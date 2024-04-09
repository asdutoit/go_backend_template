package routes

import (
	"os"

	"github.com/asdutoit/go_backend_template/graphql"
	"github.com/asdutoit/go_backend_template/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()

	if os.Getenv("ENV") == "production" {
		config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))

	// Add your routes here
	r.GET("/healthcheck", healthCheck)
	r.POST("/signup", signUp)
	r.POST("/login", login)
	r.POST("/events", middlewares.Authenticate, createEvent)
	r.GET("/deployments", getDeployments)
	r.GET("/deployment", GetDeploymentByQuery)
	r.GET("/auth/google", handleGoogleAuth)
	r.GET("/auth/google/callback", handleGoogleCallback)
	r.GET("/metrics", getMetrics)
	r.Use(MetricsMiddleware())
	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.GET("/events", getEvents)
	authenticated.GET("/events/:id", getEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.GET("/users", getUsers)
	authenticated.GET("/user", getUser)
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

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of incoming HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	// Register the metrics.
	prometheus.MustRegister(httpRequestsTotal)
}

func IncrementRequestsCounter(method, endpoint string) {
	httpRequestsTotal.With(prometheus.Labels{"method": method, "endpoint": endpoint}).Inc()
}

func getMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		IncrementRequestsCounter(c.Request.Method, c.FullPath())
		c.Next()
	}
}
