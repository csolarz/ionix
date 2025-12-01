package controller

import (
	"os"

	"github.com/csolarz/ionix/internal/infra"
	"github.com/csolarz/ionix/internal/usecase"
)

type dependencies struct {
	Auth *AuthController
	Task *TaskController
}

func RegisterDependencies() dependencies {
	db, err := infra.NewDBGorm(os.Getenv("CONNECTION_STRING_DB"))
	if err != nil {
		panic(err)
	}

	authSvc := usecase.NewAuthService(db)
	authCtrl := NewAuthController(authSvc)

	taskSvc := usecase.NewTaskService(db)
	taskCtrl := NewTaskController(taskSvc)

	return dependencies{
		Auth: authCtrl,
		Task: taskCtrl,
	}
}
