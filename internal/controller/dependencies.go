package controller

import (
	"os"

	"github.com/csolarz/ionix/internal/infra"
	"github.com/csolarz/ionix/internal/usecase"
)

type dependencies struct {
	*AuthController
	*TaskController
}

func registerDependencies() dependencies {
	db, err := infra.NewDBGorm(os.Getenv("CONECTION_STRING_DB"))
	if err != nil {
		panic(err)
	}

	taskSvc := usecase.NewTaskService(db)
	taskCtrl := NewTaskController(taskSvc)

	return dependencies{
		AuthController: NewAuthController(nil),
		TaskController: taskCtrl,
	}
}
