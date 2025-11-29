package usecase

import (
	"context"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/infra"
)

// Define los casos de uso relacionados con las tareas, sirve para desacoplar la lógica de negocio
// de la implementación concreta del repositorio de datos.
// Se usa mockery para generar mocks de esta interfaz para pruebas unitarias.
//
//go:generate mockery --name=TaskUsecase --output=./mock --outpkg=mock --case=snake
type TaskUsecase interface {
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
	GetAll(ctx context.Context) ([]domain.Task, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskService struct {
	repo infra.DBRepository
}

func NewTaskService(repo infra.DBRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	err := s.repo.Create(ctx, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	var task *domain.Task
	err := s.repo.GetByID(ctx, id, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	data := []domain.Task{}
	err := s.repo.GetAll(ctx, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *TaskService) Update(ctx context.Context, task domain.Task) error {
	return s.repo.Update(ctx, &task)
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
