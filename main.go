package main

import (
	"log"

	"github.com/csolarz/ionix/internal/controller"
)

func main() {
	router := controller.SetupRouter()

	// escucha en 0.0.0.0:8080 por defecto
	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
