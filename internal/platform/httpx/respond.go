package httpx

import (
	"log/slog"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
	"github.com/Hello256World/cloud-native-shop/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  []apierrors.FieldError `json:"fields,omitempty"`
}

func statusAndLevel(k apierrors.Kind) (int, slog.Level) {
	switch k {
	case apierrors.KindInvalidInput:
		return http.StatusBadRequest, slog.LevelWarn
	case apierrors.KindNotFound:
		return http.StatusNotFound, slog.LevelInfo
	case apierrors.KindConflict:
		return http.StatusConflict, slog.LevelWarn
	case apierrors.KindUnauthorized:
		return http.StatusUnauthorized, slog.LevelInfo
	case apierrors.KindForbidden:
		return http.StatusForbidden, slog.LevelWarn
	case apierrors.KindUnimplemented:
		return http.StatusNotImplemented, slog.LevelWarn
	default:
		return http.StatusInternalServerError, slog.LevelError
	}
}

func Respond(c *gin.Context, err error) {
	ctx := c.Request.Context()
	apiErr := apierrors.FromError(err)
	status, level := statusAndLevel(apiErr.Kind())

	logger.FromContext(ctx).Log(ctx, level, "request failed",
		slog.String("code", apiErr.Code()),
		slog.Int("status", status),
		slog.Any("error", err),
	)

	c.AbortWithStatusJSON(status, ErrorBody{
		Error: ErrorDetail{Code: apiErr.Code(), Message: apiErr.Message(), Fields: apiErr.Details()},
	})
}
