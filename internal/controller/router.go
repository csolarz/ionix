package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	dependencies := registerDependencies()

	// Auth routes
	router.POST("/login", dependencies.Auth.Login)
	router.POST("/logout", dependencies.Auth.Logout)

	// User management routes
	router.POST("/users", dependencies.Auth.Register)
	router.PUT("/users", dependencies.Auth.UpdatePassword)
	router.DELETE("/users/:id", dependencies.Auth.Delete)

	// Task routes
	router.POST("/tasks", dependencies.Task.CreateTask)
	router.GET("/tasks/:id", dependencies.Task.GetTaskByID)
	router.GET("/tasks", dependencies.Task.GetAllTasks)
	router.PUT("/tasks/:id", dependencies.Task.UpdateTask)
	router.DELETE("/tasks/:id", dependencies.Task.DeleteTask)

	// Health check route
	router.GET("/ping", pingController)

	return router
}

func pingController(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
