package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/gin-gonic/gin"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func newRouter(tm *auth.TokenManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/me", Auth(tm), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/admin", Auth(tm), RequireRole("admin"), func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func do(r http.Handler, path, header string) int {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func TestAuthAndRoles(t *testing.T) {
	tm := auth.NewTokenManager(testSecret, "test", 15*time.Minute)
	r := newRouter(tm)

	customer, _, err := tm.Issue(1, []string{"customer"})
	if err != nil {
		t.Fatal(err)
	}
	admin, _, err := tm.Issue(2, []string{"admin"})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name, path, header string
		want               int
	}{
		{"no header", "/me", "", http.StatusUnauthorized},
		{"wrong scheme", "/me", "Basic abc", http.StatusUnauthorized},
		{"garbage token", "/me", "Bearer not.a.jwt", http.StatusUnauthorized},
		{"valid customer", "/me", "Bearer " + customer, http.StatusOK},
		{"customer on admin", "/admin", "Bearer " + customer, http.StatusForbidden},
		{"admin on admin", "/admin", "Bearer " + admin, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := do(r, tc.path, tc.header); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}