package controller

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/csolarz/ionix/internal/domain"
	"github.com/csolarz/ionix/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	usecase usecase.AuthUsecase
}

// es la clave usada para identificar al usuario en el token JWT
var identityKey = "id"

func NewAuthController(usecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		usecase: usecase,
	}
}

func (api *AuthController) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

func (api *AuthController) UpdatePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "Not implemented")
}

// devuelve las claims JWT para el usuario autenticado
func (api *AuthController) PayloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(*domain.User); ok {
		return jwt.MapClaims{
			identityKey: user.ID,
			"username":  user.Username,
			"role":      user.Role,
		}
	}
	return jwt.MapClaims{}
}

func (api *AuthController) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &domain.User{
		ID:       int64(claims[identityKey].(float64)),
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	}
}

func (api *AuthController) Authenticator(c *gin.Context) (interface{}, error) {
	var login domain.User
	if err := c.ShouldBindJSON(&login); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	if login.Username == "" || login.Password == "" {
		return nil, jwt.ErrMissingLoginValues
	}

	user := &domain.User{
		Username: login.Username,
		Password: login.Password,
	}

	//user, err := api.usecase.Login(c.Request.Context(), login.Username, login.Password)
	//if err != nil {
	//	return nil, jwt.ErrFailedAuthentication
	//}

	return user, nil
}
