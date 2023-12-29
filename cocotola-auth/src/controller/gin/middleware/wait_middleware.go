package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewWaitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(1 * time.Second)
	}
}
