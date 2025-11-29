package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/csolarz/ionix/internal/domain"
	mockusecase "github.com/csolarz/ionix/internal/usecase/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouterWithMock(m *mockusecase.TaskUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	tc := NewTaskController(m)
	r := gin.New()
	r.POST("/tasks", tc.CreateTask)
	r.GET("/tasks/:id", tc.GetTaskByID)
	r.GET("/tasks", tc.GetAllTasks)
	r.PUT("/tasks/:id", tc.UpdateTask)
	r.DELETE("/tasks/:id", tc.DeleteTask)
	return r
}

func TestCreateTask_Success_WithGeneratedMock(t *testing.T) {
	m := mockusecase.NewTaskUsecase(t)
	m.On("Create", mock.Anything, mock.Anything).Return(&domain.Task{ID: 1, Title: "Test Task"}, nil)

	router := setupRouterWithMock(m)

	body := []byte(`{"title":"Test Task","due_date": "2025-12-12T00:00:00Z"}`)
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

	body := []byte(`{"title":"Test Task","due_date": "2025-12-12T00:00:00Z"}`)
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
