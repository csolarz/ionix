package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/csolarz/ionix/internal/domain"
	mockrepo "github.com/csolarz/ionix/internal/infra/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateCredentials_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	// Expect Validate to be called and succeed
	m.On("Validate", mock.Anything, mock.Anything).Return(nil)

	svc := NewAuthService(m)
	u := &domain.User{Username: "alice", Password: "secret"}

	err := svc.ValidateCredentials(context.Background(), u)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestValidateCredentials_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	m.On("Validate", mock.Anything, mock.Anything).Return(errors.New("invalid credentials"))

	svc := NewAuthService(m)
	u := &domain.User{Username: "bob", Password: "bad"}

	err := svc.ValidateCredentials(context.Background(), u)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid credentials")

	m.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	// Verify that Create is called with a *domain.User that has the expected fields
	m.On("Create", mock.Anything, mock.MatchedBy(func(v interface{}) bool {
		user, ok := v.(*domain.User)
		if !ok {
			return false
		}
		return user.Username == "newuser" && user.Password == "pass123" && user.Role == "user"
	})).Return(nil)

	svc := NewAuthService(m)
	err := svc.Register(context.Background(), "newuser", "pass123")
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestRegister_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	m.On("Create", mock.Anything, mock.Anything).Return(errors.New("create failed"))

	svc := NewAuthService(m)
	err := svc.Register(context.Background(), "newuser", "pass123")
	assert.Error(t, err)
	assert.EqualError(t, err, "create failed")

	m.AssertExpectations(t)
}

func TestLogout_Noop(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewAuthService(m)

	// Logout currently is a noop and should return nil
	err := svc.Logout(context.Background(), "some-user-id")
	assert.NoError(t, err)
	// No expectations to assert (no repo calls needed)
}

func TestUpdatePassword_Noop(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewAuthService(m)

	// UpdatePassword currently returns nil
	err := svc.UpdatePassword(context.Background(), "user-id", "newpass")
	assert.NoError(t, err)
}
