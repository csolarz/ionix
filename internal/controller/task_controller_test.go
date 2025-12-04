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

func TestGetTaskByID_Success_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("GetByID", mock.Anything, int64(1)).Return(&domain.Task{ID: 1, Title: "found"}, nil)

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "found", resp.Title)

	m.AssertExpectations(t)
}

func TestGetTaskByID_InvalidID_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTaskByID_NotFound_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("GetByID", mock.Anything, int64(999)).Return((*domain.Task)(nil), nil)

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}

func TestGetTaskByID_UsecaseError_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("GetByID", mock.Anything, int64(1)).Return((*domain.Task)(nil), errors.New("repo error"))

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

// ======= New tests for GetAllTasks, UpdateTask, DeleteTask =======

func TestGetAllTasks_Success(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	list := []domain.Task{
		{ID: 1, Title: "A"},
		{ID: 2, Title: "B"},
	}
	m.On("GetAll", mock.Anything).Return(list, nil)

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, int64(1), resp[0].ID)

	m.AssertExpectations(t)
}

func TestGetAllTasks_Error(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("GetAll", mock.Anything).Return(([]domain.Task)(nil), errors.New("db"))

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestUpdateTask_Success(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	// Expect Update to be called with a task that has ID == 1 and Title == "Updated"
	m.On("Update", mock.Anything, mock.MatchedBy(func(t domain.Task) bool {
		return t.ID == 1 && t.Title == "Updated"
	})).Return(nil)

	router := setupRouterWithMock(m)

	body := []byte(`{"id":1,"title":"Updated", "due_date":"2025-12-12T00:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// body is JSON string, check it includes expected message
	assert.Contains(t, w.Body.String(), "Task updated successfully")

	m.AssertExpectations(t)
}

func TestUpdateTask_InvalidID(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	router := setupRouterWithMock(m)

	body := []byte(`{"title":"Updated"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTask_InvalidJSON(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTask_UsecaseError(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed"))

	router := setupRouterWithMock(m)

	body := []byte(`{"id":1,"title":"Updated", "due_date":"2025-12-12T00:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteTask_Success(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Delete", mock.Anything, 1).Return(nil)

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestDeleteTask_InvalidID(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTask_UsecaseError(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Delete", mock.Anything, 1).Return(errors.New("delete failed"))

	router := setupRouterWithMock(m)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}
