package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInit_CreatesLogger(t *testing.T) {
	// Reset Log to nil before test
	Log = nil

	Init()

	assert.NotNil(t, Log)
	assert.IsType(t, (*zap.SugaredLogger)(nil), Log)
}

func TestInit_CalledMultipleTimes(t *testing.T) {
	// Reset Log
	Log = nil

	Init()
	firstLog := Log
	assert.NotNil(t, firstLog)

	// Call Init again
	Init()
	secondLog := Log
	assert.NotNil(t, secondLog)

	// Both should be valid logger instances (may be different objects)
	assert.IsType(t, (*zap.SugaredLogger)(nil), secondLog)
}

func TestSync_WithNilLogger(t *testing.T) {
	// Reset Log to nil
	Log = nil

	// Should not panic
	assert.NotPanics(t, func() {
		Sync()
	})
}

func TestSync_WithValidLogger(t *testing.T) {
	// Initialize logger first
	Log = nil
	Init()

	// Should not panic when calling Sync on valid logger
	assert.NotPanics(t, func() {
		Sync()
	})
}

func TestLogWithInfo(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log without panicking
	assert.NotPanics(t, func() {
		Log.Info("test info message")
	})
}

func TestLogWithError(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log errors without panicking
	assert.NotPanics(t, func() {
		Log.Error("test error message")
	})
}

func TestLogWithWarn(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log warnings without panicking
	assert.NotPanics(t, func() {
		Log.Warn("test warn message")
	})
}

func TestLogWithDebug(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log debug without panicking
	assert.NotPanics(t, func() {
		Log.Debug("test debug message")
	})
}

func TestLogWithFields(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log with fields without panicking
	assert.NotPanics(t, func() {
		Log.Infow("test message with fields",
			"user_id", 42,
			"action", "create",
			"status", "success",
		)
	})
}

func TestLogErrorWithFields(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Test that we can log errors with fields without panicking
	assert.NotPanics(t, func() {
		Log.Errorw("error occurred",
			"error_code", 500,
			"service", "auth_service",
		)
	})
}

func TestSync_TwiceDoesNotPanic(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Call Sync twice - should not panic
	assert.NotPanics(t, func() {
		Sync()
		Sync()
	})
}

func TestLoggerIsProduction(t *testing.T) {
	// Initialize logger
	Log = nil
	Init()

	// Production logger should be configured
	assert.NotNil(t, Log)

	// Verify it's a SugaredLogger
	assert.IsType(t, (*zap.SugaredLogger)(nil), Log)
}

func TestInitErrorHandling(t *testing.T) {
	// Init ignores errors from zap.NewProduction
	// This test ensures Init doesn't panic even if there's an issue
	Log = nil

	assert.NotPanics(t, func() {
		Init()
	})

	assert.NotNil(t, Log)
}
