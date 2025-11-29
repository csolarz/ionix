package controller

import (
	"net/http"

	"github.com/csolarz/ionix/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	usecase usecase.AuthUsecase
}

func NewAuthController(usecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		usecase: usecase,
	}
}

func (api *AuthController) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) UpdatePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}
