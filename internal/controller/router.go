package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	dependencies := registerDependencies()

	// Auth routes
	router.POST("/login", dependencies.AuthController.Login)
	router.POST("/logout", dependencies.AuthController.Logout)

	// User management routes
	router.POST("/users", dependencies.AuthController.Register)
	router.PUT("/users", dependencies.AuthController.UpdatePassword)
	router.DELETE("/users/:id", dependencies.AuthController.Delete)

	// Task routes
	router.POST("/tasks", dependencies.TaskController.CreateTask)
	router.GET("/tasks/:id", dependencies.TaskController.GetTaskByID)
	router.GET("/tasks", dependencies.TaskController.GetAllTasks)
	router.PUT("/tasks/:id", dependencies.TaskController.UpdateTask)
	router.DELETE("/tasks/:id", dependencies.TaskController.DeleteTask)

	// Health check route
	router.GET("/ping", pingController)

	return router
}

func pingController(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
