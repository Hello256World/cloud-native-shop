package http

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return f.Name
			}
			return name
		})
	}
}

func bind(c *gin.Context, req any) bool {
	err := c.ShouldBind(req)
	if err == nil {
		return true
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		details := make([]apierrors.FieldError, len(verrs))
		for i, fe := range verrs {
			details[i] = apierrors.FieldError{
				Field: fe.Field(),
				Code:  fe.Tag(),
				Param: fe.Param(),
			}
		}
		httpx.Respond(c, apierrors.Wrap(apierrors.NewValidationError(details), err))
		return false
	}

	httpx.Respond(c, apierrors.Wrap(
		apierrors.NewBadRequestError("invalid_request", "invalid request body"), err))
	return false
}

func currentUserID(c *gin.Context) (int64, bool) {
	p, ok := auth.FromContext(c.Request.Context())

	if !ok {
		httpx.Respond(c, apierrors.NewUnauthorizedError("missing_token", "authentication required"))
		return 0, false
	}

	return p.UserID, true
}

func pathID(c *gin.Context) (int64, bool) {
	raw := c.Param("id")

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpx.Respond(c, apierrors.Wrap(
			apierrors.NewBadRequestError("invalid_id", "invalid id"), fmt.Errorf("invalid path id %q", raw)))

		return 0, false
	}

	return id, true
}
