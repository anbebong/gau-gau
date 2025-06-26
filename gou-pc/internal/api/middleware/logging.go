package middleware

import (
	"gou-pc/internal/logutil"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logutil.Info("[API] %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
