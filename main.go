package main

import (
	"log"
	"os"

	"github.com/asdutoit/gotraining/section11/db"
	"github.com/asdutoit/gotraining/section11/routes"
)

func main() {
	checkEnv()
	db.InitDB()
	server := routes.SetupRouter()
	server.ForwardedByClientIP = true
	server.SetTrustedProxies([]string{"127.0.0.1"})

	server.Run(":8080")
}

func checkEnv() {
	envVars := []string{"AWS_REGION", "AWS_S3_BUCKET", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DB_NAME", "DB_PASSWORD", "DB_USER", "DB_HOST", "DB_PORT"}
	missingVars := []string{}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}
	if len(missingVars) > 0 {
		log.Fatalf("Missing environment variables: %v", missingVars)
	}
}
