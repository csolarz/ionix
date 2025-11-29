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

func TestCreate_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	in := &domain.Task{Title: "new"}

	m.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// opcional: podrías modificar el valor pasado si tu repo lo hace
	})

	got, err := svc.Create(context.Background(), in)
	assert.NoError(t, err)
	assert.Equal(t, in, got)

	m.AssertExpectations(t)
}

func TestCreate_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("Create", mock.Anything, mock.Anything).Return(errors.New("create failed"))

	_, err := svc.Create(context.Background(), &domain.Task{})
	assert.Error(t, err)

	m.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	id := int64(42)
	m.On("GetByID", mock.Anything, id, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// args.Get(2) is **domain.Task; asignamos el resultado esperado
		if ptr, ok := args.Get(2).(**domain.Task); ok {
			*ptr = &domain.Task{ID: id, Title: "found"}
		}
	})

	got, err := svc.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "found", got.Title)

	m.AssertExpectations(t)
}

func TestGetByID_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found"))

	_, err := svc.GetByID(context.Background(), 1)
	assert.Error(t, err)

	m.AssertExpectations(t)
}

func TestGetAll_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("GetAll", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if ptr, ok := args.Get(1).(*[]domain.Task); ok {
			*ptr = []domain.Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}}
		}
	})

	list, err := svc.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(1), list[0].ID)

	m.AssertExpectations(t)
}

func TestGetAll_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("GetAll", mock.Anything, mock.Anything).Return(errors.New("get all failed"))

	_, err := svc.GetAll(context.Background())
	assert.Error(t, err)

	m.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), domain.Task{ID: 5, Title: "x"})
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestUpdate_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed"))

	err := svc.Update(context.Background(), domain.Task{ID: 5})
	assert.Error(t, err)

	m.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	// Delete signature in mock is (ctx, data interface{})
	m.On("Delete", mock.Anything, 7).Return(nil)

	err := svc.Delete(context.Background(), 7)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	m := mockrepo.NewDBRepository(t)
	svc := NewTaskService(m)

	m.On("Delete", mock.Anything, 7).Return(errors.New("delete failed"))

	err := svc.Delete(context.Background(), 7)
	assert.Error(t, err)

	m.AssertExpectations(t)
}
