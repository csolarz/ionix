package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Aquí se captura el panic, se puede registrar el error y devolver una respuesta adecuada
				log.Printf("PANIC: %v\nStack: %s", err, debug.Stack())

				c.JSON(500, gin.H{
					"error": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
