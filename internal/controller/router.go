package controller

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/csolarz/ionix/internal/domain/constants"
	"github.com/csolarz/ionix/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dependencies dependencies) *gin.Engine {
	router := gin.Default()

	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "example zone",
		Key:         []byte("secret_here"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",

		Authenticator:   dependencies.Auth.Authenticator,
		PayloadFunc:     dependencies.Auth.PayloadFunc,
		IdentityHandler: dependencies.Auth.IdentityHandler,

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
	})

	// Health check
	router.GET("/ping", pingController)

	// Login
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/logout", dependencies.Auth.Logout)

	// Protected routes
	router.Use(authMiddleware.MiddlewareFunc())

	api := router.Group("/api")

	// User routes (only for admin)
	api.PUT("/users", dependencies.Auth.UpdatePassword)
	api.POST("/users", middleware.RoleRequired(constants.RoleAdmin), dependencies.Auth.Register)
	api.GET("/users", middleware.RoleRequired(constants.RoleAdmin), dependencies.Auth.Register)

	// Task routes
	api.POST("/tasks", dependencies.Task.CreateTask)
	api.GET("/tasks/:id", dependencies.Task.GetTaskByID)
	api.GET("/tasks", dependencies.Task.GetAllTasks)
	api.PUT("/tasks/:id", dependencies.Task.UpdateTask)
	api.DELETE("/tasks/:id", dependencies.Task.DeleteTask)

	return router
}

func pingController(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
