package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(context.Context, usecase.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, phone, password string) (*usecase.TokenPair, error)
	Refresh(context.Context, string) (*usecase.TokenPair, error)
	Logout(context.Context, string) error
}

type Handler struct {
	svc authService
}

func NewHandler(svc authService) *Handler { return &Handler{svc: svc} }

func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBind(req); err != nil {
		httpx.Respond(c, apierrors.Wrap(
			apierrors.NewBadRequestError("invalid_request", "invalid request body"), err))

		return false
	}

	return true
}

// Register godoc
// @Summary     Register a new customer
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body RegisterRequest true "Registration data"
// @Success     201 {object} UserResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if !bind(c, &req) {
		return
	}

	u, err := h.svc.Register(c.Request.Context(), usecase.RegisterInput{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(u))
}

// Login godoc
// @Summary     Login with phone and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body LoginRequest true "Credentials"
// @Success     200 {object} TokenResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Router      /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if !bind(c, &req) {
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Phone, req.Password)

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}
	c.JSON(http.StatusAccepted, toTokenResponse(token))
}

// Refresh godoc
// @Summary     Rotate the refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body RefreshRequest true "Refresh token"
// @Success     200 {object} TokenResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Router      /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if !bind(c, &req) {
		return
	}

	token, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toTokenResponse(token))
}

// Logout godoc
// @Summary     Revoke the refresh token
// @Tags        auth
// @Accept      json
// @Param       request body RefreshRequest true "Refresh token"
// @Success     204
// @Failure     400 {object} httpx.ErrorBody
// @Router      /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest

	if !bind(c, &req) {
		return
	}

	if err := h.svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
