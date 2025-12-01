package usecase

import "context"

type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, userID string) error
	Register(ctx context.Context, username, password string) error
	UpdatePassword(ctx context.Context, userID, newPassword string) error
}
