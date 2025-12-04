package main

import (
	"log"
	"os"

	"github.com/csolarz/ionix/internal/controller"
	"github.com/csolarz/ionix/internal/infra/logger"
)

const defaultPort = "8080"

func main() {
	logger.Init()
	defer logger.Sync()

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
