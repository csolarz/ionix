package middleware

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		role := claims["role"].(string)

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "no tienes permiso para acceder a esta ruta",
		})
		c.Abort()
	}
}
