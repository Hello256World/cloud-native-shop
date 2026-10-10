package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Hello256World/cloud-native-shop/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.FromContext(c.Request.Context()).Error("panic recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "internal_error", "message": "internal server error"},
				})
			}
		}()

		c.Next()
	}
}
