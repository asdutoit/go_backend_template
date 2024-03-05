package main

import (
	"github.com/asdutoit/gotraining/section11/db"
	"github.com/asdutoit/gotraining/section11/routes"
)

func main() {
	db.InitDB()
	server := routes.SetupRouter()
	server.ForwardedByClientIP = true
	server.SetTrustedProxies([]string{"127.0.0.1"})

	server.Run(":8080")
}
