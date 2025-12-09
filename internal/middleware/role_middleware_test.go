package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// helper: create a router that injects JWT claims.role and registers a test route with RoleRequired
func routerWithRoleAndRequired(role string, requiredRoles ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// middleware to inject JWT_PAYLOAD claims into context
	r.Use(func(c *gin.Context) {
		c.Set("JWT_PAYLOAD", jwt.MapClaims{
			"role": role,
		})
		c.Next()
	})

	// register a GET /test that uses RoleRequired and writes "ok" when permitted
	r.GET("/test", RoleRequired(requiredRoles...), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	})

	return r
}

func TestRoleRequired_AllowedRole(t *testing.T) {
	router := routerWithRoleAndRequired("admin", "admin")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["result"])
}

func TestRoleRequired_DeniedRole(t *testing.T) {
	router := routerWithRoleAndRequired("user", "admin")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "no tienes permiso")
}

func TestRoleRequired_MultipleAllowedRoles(t *testing.T) {
	router := routerWithRoleAndRequired("moderator", "admin", "moderator")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleRequired_EmptyRequiredRoles_DeniesAll(t *testing.T) {
	// If no roles are passed to RoleRequired, current implementation will deny access
	router := routerWithRoleAndRequired("admin") // no required roles provided to middleware
	// register route with no required roles
	// Note: the helper registers RoleRequired() with no roles which will deny all

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
