package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type addressService interface {
	AddAddress(c context.Context, userID int64, in usecase.AddressInput) (*domain.Address, error)
	ListAddress(c context.Context, userID int64) ([]domain.Address, error)
	UpdateAddress(c context.Context, userID, id int64, in usecase.AddressInput) (*domain.Address, error)
	DeleteAddress(c context.Context, userID, id int64) error
}

type AddressHandler struct {
	svc addressService
}

func NewAddressHandler(svc addressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func toAddressInput(r AddressRequest) usecase.AddressInput {
	return usecase.AddressInput{
		ReceiverName: r.ReceiverName,
		Phone:        r.Phone,
		Address:      r.Address,
		PostalCode:   r.PostalCode,
		No:           r.No,
		Unit:         r.Unit,
	}
}

// Create godoc
// @Summary     Add an address
// @Tags        me
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body AddressRequest true "Address"
// @Success     201 {object} AddressResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Router      /me/addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}

	var req AddressRequest
	if !bind(c, &req) {
		return
	}

	addr, err := h.svc.AddAddress(c.Request.Context(), id, toAddressInput(req))

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toAddressResponse(addr))
}

// List godoc
// @Summary     List my addresses
// @Tags        me
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array}  AddressResponse
// @Failure     401 {object} httpx.ErrorBody
// @Router      /me/addresses [get]
func (h *AddressHandler) List(c *gin.Context) {
	id, ok := currentUserID(c)

	if !ok {
		return
	}

	addrs, err := h.svc.ListAddress(c.Request.Context(), id)

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toAddressListResponse(addrs))
}

// Update godoc
// @Summary     Replace an address
// @Tags        me
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Address ID"
// @Param       request body AddressRequest true "Address"
// @Success     200 {object} AddressResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /me/addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)

	if !ok {
		return
	}

	id, ok := pathID(c)

	if !ok {
		return
	}

	var req AddressRequest

	if !bind(c, &req) {
		return
	}

	addr, err := h.svc.UpdateAddress(c.Request.Context(), userID, id, toAddressInput(req))

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toAddressResponse(addr))
}

// Delete godoc
// @Summary     Delete an address
// @Tags        me
// @Security    BearerAuth
// @Param       id path int true "Address ID"
// @Success     204
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /me/addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	id, ok := pathID(c)
	if !ok {
		return
	}

	if err := h.svc.DeleteAddress(c.Request.Context(), userID, id); err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
