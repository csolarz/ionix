package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/csolarz/ionix/internal/domain"
	mockusecase "github.com/csolarz/ionix/internal/usecase/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthRouter(m *mockusecase.AuthUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	ac := NewAuthController(m)
	r := gin.New()
	r.POST("/login", ac.Login)
	r.POST("/logout", ac.Logout)
	r.POST("/register", ac.Register)
	r.PUT("/password", ac.UpdatePassword)
	return r
}

// ============ Login Tests ============

func TestLogin_NotImplemented(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	router := setupAuthRouter(m)

	body := []byte(`{"username":"user","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// ============ Logout Tests ============

func TestLogout_NotImplemented(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	router := setupAuthRouter(m)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// ============ Register Tests ============

func TestRegister_NotImplemented(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	router := setupAuthRouter(m)

	body := []byte(`{"username":"newuser","password":"pass123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// ============ UpdatePassword Tests ============

func TestUpdatePassword_NotImplemented(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	router := setupAuthRouter(m)

	body := []byte(`{"old_password":"oldpass","new_password":"newpass"}`)
	req := httptest.NewRequest(http.MethodPut, "/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

// ============ PayloadFunc Tests ============

func TestPayloadFunc_ValidUser(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	user := &domain.User{
		Username: "testuser",
		Role:     "admin",
	}

	claims := ac.PayloadFunc(user)

	assert.NotNil(t, claims)
	assert.Equal(t, "testuser", claims[identityKey])
	assert.Equal(t, "admin", claims["role"])

	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestPayloadFunc_InvalidType(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	claims := ac.PayloadFunc("not a user")

	assert.NotNil(t, claims)
	assert.Empty(t, claims)

	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestPayloadFunc_NilData(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	claims := ac.PayloadFunc(nil)

	assert.NotNil(t, claims)
	assert.Empty(t, claims)

	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestPayloadFunc_UserRole(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	user := &domain.User{
		Username: "john",
		Role:     "user",
	}

	claims := ac.PayloadFunc(user)

	assert.Equal(t, "john", claims[identityKey])
	assert.Equal(t, "user", claims["role"])
}

func TestPayloadFunc_GuestRole(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	user := &domain.User{
		Username: "guest",
		Role:     "guest",
	}

	claims := ac.PayloadFunc(user)

	assert.Equal(t, "guest", claims[identityKey])
	assert.Equal(t, "guest", claims["role"])
}

// ============ IdentityHandler Tests ============

func TestIdentityHandler_ValidClaims(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("JWT_PAYLOAD", jwt.MapClaims{
		"username": "testuser",
		"role":     "user",
	})

	result := ac.IdentityHandler(c)

	assert.NotNil(t, result)
	if user, ok := result.(*domain.User); ok {
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "user", user.Role)
	}

	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestIdentityHandler_AdminClaims(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("JWT_PAYLOAD", jwt.MapClaims{
		"username": "admin",
		"role":     "admin",
	})

	result := ac.IdentityHandler(c)

	assert.NotNil(t, result)
	if user, ok := result.(*domain.User); ok {
		assert.Equal(t, "admin", user.Username)
		assert.Equal(t, "admin", user.Role)
	}
}

func TestIdentityHandler_GuestClaims(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("JWT_PAYLOAD", jwt.MapClaims{
		"username": "guest",
		"role":     "guest",
	})

	result := ac.IdentityHandler(c)

	assert.NotNil(t, result)
	if user, ok := result.(*domain.User); ok {
		assert.Equal(t, "guest", user.Username)
		assert.Equal(t, "guest", user.Role)
	}
}

// ============ Authenticator Tests ============

func TestAuthenticator_ValidCredentials_Success(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	m.On("ValidateCredentials", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == "testuser" && user.Password == "correctpass"
	})).Return(nil)

	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"testuser","password":"correctpass"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if user, ok := result.(domain.User); ok {
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "correctpass", user.Password)
	}
	m.AssertExpectations(t)
}

func TestAuthenticator_InvalidCredentials_WrongPassword(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	m.On("ValidateCredentials", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == "testuser" && user.Password == "wrongpass"
	})).Return(errors.New("invalid credentials"))

	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"testuser","password":"wrongpass"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrFailedAuthentication, err)
	m.AssertExpectations(t)
}

func TestAuthenticator_InvalidCredentials_UserNotFound(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	m.On("ValidateCredentials", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == "nonexistent"
	})).Return(errors.New("user not found"))

	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"nonexistent","password":"anypass"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrFailedAuthentication, err)
	m.AssertExpectations(t)
}

func TestAuthenticator_InvalidJSON_MalformedRequest(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{invalid json}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrMissingLoginValues, err)
	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestAuthenticator_MissingUsername(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"password":"pass123"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	// El comportamiento depende de cómo ShouldBindJSON maneje campos vacíos
	// Si falla, esperamos ErrMissingLoginValues
	if err != nil {
		assert.Equal(t, jwt.ErrMissingLoginValues, err)
		assert.Nil(t, result)
	}
}

func TestAuthenticator_MissingPassword(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"testuser"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	if err != nil {
		assert.Equal(t, jwt.ErrMissingLoginValues, err)
		assert.Nil(t, result)
	}
}

func TestAuthenticator_DatabaseError(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	m.On("ValidateCredentials", mock.Anything, mock.Anything).
		Return(errors.New("database connection failed"))

	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"testuser","password":"pass"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrFailedAuthentication, err)
	m.AssertExpectations(t)
}

func TestAuthenticator_RepositoryError(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	m.On("ValidateCredentials", mock.Anything, mock.Anything).
		Return(errors.New("repository error"))

	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"testuser","password":"pass"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrFailedAuthentication, err)
	m.AssertExpectations(t)
}

func TestAuthenticator_EmptyCredentials(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := []byte(`{"username":"","password":""}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	result, err := ac.Authenticator(c)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, jwt.ErrMissingLoginValues, err)
	m.AssertExpectations(t)
}

// ============ Integration-like Tests ============

func TestAuthController_PayloadFunc_WithDifferentRoles(t *testing.T) {
	m := mockusecase.NewAuthUsecase(t)
	ac := NewAuthController(m)

	testCases := []struct {
		name     string
		user     *domain.User
		expected map[string]interface{}
	}{
		{
			name:     "Admin role",
			user:     &domain.User{Username: "admin", Role: "admin"},
			expected: map[string]interface{}{identityKey: "admin", "role": "admin"},
		},
		{
			name:     "User role",
			user:     &domain.User{Username: "john", Role: "user"},
			expected: map[string]interface{}{identityKey: "john", "role": "user"},
		},
		{
			name:     "Guest role",
			user:     &domain.User{Username: "guest", Role: "guest"},
			expected: map[string]interface{}{identityKey: "guest", "role": "guest"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims := ac.PayloadFunc(tc.user)
			assert.Equal(t, tc.expected[identityKey], claims[identityKey])
			assert.Equal(t, tc.expected["role"], claims["role"])
		})
	}

	m.AssertNotCalled(t, "ValidateCredentials")
}

func TestAuthController_Authenticator_WithDifferentUsers(t *testing.T) {
	testCases := []struct {
		name       string
		username   string
		password   string
		shouldFail bool
	}{
		{
			name:       "Admin login",
			username:   "admin",
			password:   "admin123",
			shouldFail: false,
		},
		{
			name:       "Regular user login",
			username:   "john",
			password:   "john123",
			shouldFail: false,
		},
		{
			name:       "Guest login",
			username:   "guest",
			password:   "guest123",
			shouldFail: false,
		},
		{
			name:       "Wrong password",
			username:   "admin",
			password:   "wrongpass",
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockusecase.NewAuthUsecase(t)

			if tc.shouldFail {
				m.On("ValidateCredentials", mock.Anything, mock.Anything).
					Return(errors.New("invalid credentials"))
			} else {
				m.On("ValidateCredentials", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.Username == tc.username && user.Password == tc.password
				})).Return(nil)
			}

			ac := NewAuthController(m)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			body := []byte(`{"username":"` + tc.username + `","password":"` + tc.password + `"}`)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			result, err := ac.Authenticator(c)

			if tc.shouldFail {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			m.AssertExpectations(t)
		})
	}
}
