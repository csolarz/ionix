package controller

import (
	"net/http"
	"strconv"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	usecase usecase.TaskUsecase
}

func NewTaskController(uc usecase.TaskUsecase) *TaskController {
	return &TaskController{
		usecase: uc,
	}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := tc.usecase.Create(c.Request.Context(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error creating task")
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := tc.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error retrieving task")
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, "Task not found")
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	// Implement get all tasks logic here
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	// Implement task update logic here
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	// Implement task deletion logic here
	c.JSON(http.StatusNotImplemented, "Not implemented")
}
