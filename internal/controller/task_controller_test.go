package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/infra/logger"
	mockusecase "github.com/csolarz/ionix/internal/usecase/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupRouterWithMock(m *mockusecase.TaskUsecase) *gin.Engine {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Initialize logger.Log as a no-op sugared logger to avoid test output/panics
	logger.Log = zap.NewNop().Sugar()

	tc := NewTaskController(m)
	r := gin.New()

	// Middleware for tests: inject an authenticated user so utils.GetUserID works
	r.Use(func(c *gin.Context) {
		c.Set("id", &domain.User{ID: 42, Role: "user"})
		c.Next()
	})

	r.POST("/tasks", tc.CreateTask)
	r.GET("/tasks/:id", tc.GetTaskByID)
	r.GET("/tasks", tc.GetAllTasks)
	r.PUT("/tasks/:id", tc.UpdateTask)
	r.DELETE("/tasks/:id", tc.DeleteTask)
	return r
}

func TestCreateTask_Success_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Create", mock.Anything, mock.Anything).Return(&domain.Task{ID: 1, Title: "Test Task", UserID: 42}, nil)

	router := setupRouterWithMock(m)

	body := []byte(`{"title":"Test Task","due_date":"2025-12-12T00:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "Test Task", resp.Title)
	assert.Equal(t, int64(42), resp.UserID)

	m.AssertExpectations(t)
}

func TestCreateTask_InvalidJSON_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTask_UsecaseError_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Create", mock.Anything, mock.Anything).Return((*domain.Task)(nil), errors.New("db error"))

	router := setupRouterWithMock(m)

	body := []byte(`{"title":"Test Task","due_date":"2025-12-12T00:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}
