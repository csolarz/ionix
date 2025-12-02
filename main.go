package main

import (
	"log"
	"os"

	"github.com/csolarz/ionix/internal/controller"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	dependencies := controller.RegisterDependencies()
	router := controller.SetupRouter(dependencies)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
