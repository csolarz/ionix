package utils

import (
	"errors"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/gin-gonic/gin"
)

const identityKey = "id"

// GetUser returns the full domain.User stored by the JWT middleware
func GetUser(c *gin.Context) *domain.User {
	u, exists := c.Get(identityKey)
	if !exists {
		return nil
	}

	user, ok := u.(*domain.User)
	if !ok {
		return nil
	}

	return user
}

// GetUserID returns the authenticated user ID
func GetUserID(c *gin.Context) (int64, error) {
	user := GetUser(c)
	if user == nil {
		return 0, errors.New("user not found")
	}
	return user.ID, nil
}

// GetUserRole returns the authenticated user role
func GetUserRole(c *gin.Context) (string, error) {
	user := GetUser(c)
	if user == nil {
		return "", errors.New("user not found")
	}
	return user.Role, nil
}
