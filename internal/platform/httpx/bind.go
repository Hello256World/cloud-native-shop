package httpx

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return f.Name
			}
			return name
		})
		_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
			return slugRe.MatchString(fl.Field().String())
		})
	}
}

func Bind(c *gin.Context, req any) bool {
	err := c.ShouldBind(req)
	if err == nil {
		return true
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		details := make([]apierrors.FieldError, len(verrs))
		for i, fe := range verrs {
			field := fe.Namespace()
			if j := strings.Index(field, "."); j >= 0 {
				field = field[j+1:]
			}
			details[i] = apierrors.FieldError{Field: field, Code: fe.Tag(), Param: fe.Param()}
		}
		Respond(c, apierrors.Wrap(apierrors.NewValidationError(details), err))
		return false
	}

	Respond(c, apierrors.Wrap(
		apierrors.NewBadRequestError("invalid_request", "invalid request body"), err))
	return false
}

func CurrentUserID(c *gin.Context) (int64, bool) {
	p, ok := auth.FromContext(c.Request.Context())

	if !ok {
		Respond(c, apierrors.NewUnauthorizedError("missing_token", "authentication required"))
		return 0, false
	}

	return p.UserID, true
}

func PathParamID(c *gin.Context, name string) (int64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		Respond(c, apierrors.Wrap(apierrors.NewBadRequestError("invalid_id", "invalid id"),
			fmt.Errorf("invalid path param %s=%q", name, raw)))
		return 0, false
	}
	return id, true
}

func PathID(c *gin.Context) (int64, bool) { return PathParamID(c, "id") }
