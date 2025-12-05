package infra

import (
	"context"
	"testing"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *DBGorm {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&domain.Task{}, &domain.User{})
	require.NoError(t, err)

	return &DBGorm{db: db}
}

// TestGetByID_Success tests successful retrieval of a task by ID
func TestGetByID_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create a test task
	task := &domain.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "A test task",
		DueDate:     "2025-12-12T00:00:00Z",
		UserID:      42,
	}
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// Retrieve it
	retrieved := &domain.Task{}
	err = repo.GetByID(ctx, 1, retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", retrieved.Title)
	assert.Equal(t, int64(42), retrieved.UserID)
}

// TestGetByID_NotFound tests ErrNotFound when task doesn't exist
func TestGetByID_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	retrieved := &domain.Task{}
	err := repo.GetByID(ctx, 999, retrieved)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)
}

// TestGetAll_Success tests retrieving all tasks
func TestGetAll_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create multiple tasks
	tasks := []domain.Task{
		{Title: "Task 1", UserID: 1},
		{Title: "Task 2", UserID: 2},
	}
	for i := range tasks {
		err := repo.Create(ctx, &tasks[i])
		require.NoError(t, err)
	}

	// Retrieve all
	var retrieved []domain.Task
	err := repo.GetAll(ctx, &retrieved)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, "Task 1", retrieved[0].Title)
	assert.Equal(t, "Task 2", retrieved[1].Title)
}

// TestCreate_Success tests successful task creation
func TestCreate_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	task := &domain.Task{
		Title:  "New Task",
		UserID: 5,
	}
	err := repo.Create(ctx, task)
	assert.NoError(t, err)
	// GORM should auto-increment ID
	assert.Greater(t, task.ID, int64(0))
}

// TestUpdate_Success tests successful task update
func TestUpdate_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create
	task := &domain.Task{Title: "Original", UserID: 1}
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// Update
	task.Title = "Updated"
	err = repo.Update(ctx, task)
	assert.NoError(t, err)

	// Verify
	retrieved := &domain.Task{}
	err = repo.GetByID(ctx, task.ID, retrieved)
	require.NoError(t, err)
	assert.Equal(t, "Updated", retrieved.Title)
}

// TestDelete_Success tests successful task deletion
func TestDelete_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create
	task := &domain.Task{Title: "ToDelete", UserID: 1}
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// Delete
	err = repo.Delete(ctx, task)
	assert.NoError(t, err)

	// Verify deletion
	retrieved := &domain.Task{}
	err = repo.GetByID(ctx, task.ID, retrieved)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)
}

// TestValidate_Success tests successful user credential validation
func TestValidate_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Username: "alice",
		Password: "secret123",
		Role:     "user",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Validate with correct credentials
	toValidate := &domain.User{Username: "alice", Password: "secret123"}
	err = repo.Validate(ctx, toValidate)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, toValidate.ID) // ID should be populated
	assert.Equal(t, "user", toValidate.Role)
}

// TestValidate_WrongPassword tests validation with incorrect password
func TestValidate_WrongPassword(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		Username: "bob",
		Password: "correctpass",
		Role:     "user",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Try to validate with wrong password
	toValidate := &domain.User{Username: "bob", Password: "wrongpass"}
	err = repo.Validate(ctx, toValidate)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)
}

// TestValidate_UserNotFound tests validation when user doesn't exist
func TestValidate_UserNotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	toValidate := &domain.User{Username: "nonexistent", Password: "anypass"}
	err := repo.Validate(ctx, toValidate)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)
}
