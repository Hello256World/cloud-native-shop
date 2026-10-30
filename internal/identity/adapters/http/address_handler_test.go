package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	platformauth "github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type fakeAddrSvc struct {
	update func(userID, id int64) (*domain.Address, error)
}

func (f fakeAddrSvc) AddAddress(_ context.Context, userID int64, in usecase.AddressInput) (*domain.Address, error) {
	return &domain.Address{ID: 1, UserID: userID, ReceiverName: in.ReceiverName, Phone: in.Phone, Address: in.Address}, nil
}
func (f fakeAddrSvc) ListAddresses(context.Context, int64) ([]domain.Address, error) { return nil, nil }
func (f fakeAddrSvc) UpdateAddress(_ context.Context, userID, id int64, _ usecase.AddressInput) (*domain.Address, error) {
	return f.update(userID, id)
}
func (f fakeAddrSvc) DeleteAddress(context.Context, int64, int64) error { return nil }

func newAddrRouter(svc addressService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { // stands in for the Auth middleware
		p := platformauth.Principal{UserID: 7, Roles: []string{"customer"}}
		c.Request = c.Request.WithContext(platformauth.WithPrincipal(c.Request.Context(), p))
		c.Next()
	})
	h := NewAddressHandler(svc)
	r.POST("/me/addresses", h.Create)
	r.PUT("/me/addresses/:id", h.Update)
	return r
}

func send(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateAddress_ValidationFields(t *testing.T) {
	r := newAddrRouter(fakeAddrSvc{})
	w := send(r, http.MethodPost, "/me/addresses",
		`{"receiver_name":"Ali","phone":"09120000000","address":"Tehran, Valiasr","postal_code":"12345678901"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d, body %s", w.Code, w.Body)
	}
	var body httpx.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != "validation_failed" || len(body.Error.Fields) != 1 {
		t.Fatalf("unexpected body: %s", w.Body)
	}
	f := body.Error.Fields[0]
	if f.Field != "postal_code" || f.Code != "len" || f.Param != "10" {
		t.Fatalf("unexpected field error: %+v", f)
	}
}

func TestCreateAddress_OK(t *testing.T) {
	r := newAddrRouter(fakeAddrSvc{})
	w := send(r, http.MethodPost, "/me/addresses",
		`{"receiver_name":"Ali","phone":"09120000000","address":"Tehran, Valiasr"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d, body %s", w.Code, w.Body)
	}
}

func TestUpdateAddress_NotOwnedIs404(t *testing.T) {
	r := newAddrRouter(fakeAddrSvc{
		update: func(_, _ int64) (*domain.Address, error) { return nil, domain.ErrAddressNotFound },
	})
	w := send(r, http.MethodPut, "/me/addresses/5",
		`{"receiver_name":"Ali","phone":"09120000000","address":"Tehran, Valiasr"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status %d, body %s", w.Code, w.Body)
	}
}

func TestUpdateAddress_BadID(t *testing.T) {
	r := newAddrRouter(fakeAddrSvc{})
	w := send(r, http.MethodPut, "/me/addresses/abc", `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d", w.Code)
	}
}
