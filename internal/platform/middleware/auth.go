package middleware

import (
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
	"github.com/gin-gonic/gin"
)

type TokenVerifier interface {
	Verify(token string) (auth.Principal, error)
}

func Auth(v TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			httpx.Respond(c, apierrors.NewUnauthorizedError("missing_token", "missing or malformed authorization header"))
			return
		}

		p, err := v.Verify(raw)

		if err != nil {
			httpx.Respond(c, apierrors.Wrap(apierrors.NewUnauthorizedError("invalid_token", "invalid or expired access token"), err))
			return
		}

		c.Request = c.Request.WithContext(auth.WithPrincipal(c.Request.Context(), p))
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := auth.FromContext(c.Request.Context())

		if !ok {
			httpx.Respond(c, apierrors.NewUnauthorizedError("missing_token", "authentication required"))
			return
		}

		if !p.HasAnyRole(roles...) {
			httpx.Respond(c, apierrors.NewForbiddenError("forbidden", "you do not have access to this resource"))
			return
		}
		c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}
