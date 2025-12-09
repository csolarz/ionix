package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware_PanicHandled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Register the recovery middleware under test
	r.Use(Recovery())

	// First handler panics; second handler would set a header if executed.
	r.GET("/panic",
		func(c *gin.Context) {
			panic("boom")
		},
		func(c *gin.Context) {
			// Should NOT be executed
			c.Writer.Header().Set("X-Next-Called", "1")
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Middleware should recover and return 500 with JSON error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", body["error"])

	// Ensure the next handler was not executed
	assert.Empty(t, w.Header().Get("X-Next-Called"))
}

func TestRecoveryMiddleware_NoPanic_PassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery())

	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", body["status"])
}
