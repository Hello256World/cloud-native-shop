package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(context.Context, usecase.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, phone, password string) (*usecase.TokenPair, error)
	Refresh(context.Context, string) (*usecase.TokenPair, error)
	Logout(context.Context, string) error
	Me(context.Context, int64) (*domain.User, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc authService) *AuthHandler { return &AuthHandler{svc: svc} }

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
func (h *AuthHandler) Register(c *gin.Context) {
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
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if !bind(c, &req) {
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Phone, req.Password)

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}
	c.JSON(http.StatusOK, toTokenResponse(token))
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
func (h *AuthHandler) Refresh(c *gin.Context) {
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
func (h *AuthHandler) Logout(c *gin.Context) {
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

// Me godoc
// @Summary     Current user profile
// @Tags        me
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} UserResponse
// @Failure     401 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	u, err := h.svc.Me(c.Request.Context(), userID)
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toUserResponse(u))
}
