package usecase

import (
	"context"

	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/infra"
)

//go:generate mockery --name=AuthUsecase --output=./mock --outpkg=mock --case=snake
type AuthUsecase interface {
	ValidateCredentials(ctx context.Context, user *domain.User) error
	Logout(ctx context.Context, userID string) error
	Register(ctx context.Context, username, password string) error
	UpdatePassword(ctx context.Context, userID, newPassword string) error
}

type AuthService struct {
	repo infra.DBRepository
}

func NewAuthService(repo infra.DBRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) ValidateCredentials(ctx context.Context, user *domain.User) error {
	return s.repo.Validate(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	// Implementar la lógica de cierre de sesión si es necesario
	return nil
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	user := &domain.User{
		Username: username,
		Password: password,
		Role:     "user", // Rol por defecto
	}
	return s.repo.Create(ctx, user)
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	return nil
}
