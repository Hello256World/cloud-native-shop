package middleware

import (
	"log/slog"
	"time"

	"github.com/Hello256World/cloud-native-shop/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path

		reqID, _ := c.Get(CtxRequestID)

		l := base.With("request-id", reqID)

		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context(), l))

		c.Next()

		status := c.Writer.Status()
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"route", c.FullPath(),
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"size", c.Writer.Size(),
		}
		
		switch {
		case status >= 500:
			l.Error("http request", attrs...)
		case status >= 400:
			l.Warn("http request", attrs...)
		default:
			l.Info("http request", attrs...)
		}
	}
}
