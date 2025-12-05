package controller

import (
	"net/http"
	"strconv"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/infra/logger"
	"github.com/csolarz/ionix/internal/usecase"
	"github.com/csolarz/ionix/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	task.UserID = userID

	result, err := tc.usecase.Create(c.Request.Context(), &task)
	if err != nil {
		logger.Log.Errorw("Task creation failed",
			zap.String("task_id", task.Title),
			zap.Int64("user_id", task.UserID),
			zap.String("origin", "task_controller.CreateTask"),
		)

		c.JSON(http.StatusInternalServerError, "Error creating task")
		return
	}

	logger.Log.Infow("Task created successfully",
		zap.Int64("task_id", result.ID),
		zap.Int64("user_id", result.UserID),
		zap.String("origin", "task_controller.CreateTask"),
	)

	c.JSON(http.StatusCreated, result)
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	task, err := tc.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error retrieving tasks")
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	_id := c.Param("id")
	_, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = tc.usecase.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error updating task")
		return
	}

	c.JSON(http.StatusOK, "Task updated successfully")
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.Atoi(_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = tc.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error deleting task")
		return
	}

	c.JSON(http.StatusOK, nil)
}
