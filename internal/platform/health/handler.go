package health

import (
	"context"
	"net/http"
	"time"

	"github.com/Hello256World/cloud-native-shop/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	checkers []Checker
	timeout  time.Duration
}

func NewHandler(timeout time.Duration, checkers ...Checker) *Handler {
	return &Handler{timeout: timeout, checkers: checkers}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/healthz", h.liveness)
	r.GET("/readyz", h.readiness)
}

func (h Handler) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *Handler) readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	results := make(map[string]string, len(h.checkers))
	ok := true

	for _, ch := range h.checkers {
		if err := ch.Check(ctx); err != nil {
			results[ch.Name()] = "down"
			ok = false
			logger.FromContext(c.Request.Context()).Warn("readiness check failed", "dependency", ch.Name(), "error", err)
			continue
		}
		results[ch.Name()] = "up"
	}

	status, code := "ready", http.StatusOK

	if !ok {
		status, code = "not-ready", http.StatusServiceUnavailable
	}

	c.JSON(code, gin.H{"status": status, "checks": results})
}
