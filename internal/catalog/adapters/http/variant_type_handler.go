package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type variantTypeService interface {
	CreateVariantType(ctx context.Context, in usecase.VariantTypeInput) (*domain.VariantType, error)
	ListVariantTypes(ctx context.Context) ([]domain.VariantType, error)
	CreateVariantValue(ctx context.Context, typeID int64, in usecase.VariantValueInput) (*domain.VariantValue, error)
}

type VariantTypeHandler struct {
	svc variantTypeService
}

func NewVariantTypeHandler(svc variantTypeService) *VariantTypeHandler {
	return &VariantTypeHandler{
		svc: svc,
	}
}

// List godoc
// @Summary     List variant types with their values (admin)
// @Tags        admin-catalog
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} VariantTypeResponse
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Router      /admin/variant-types [get]
func (h *VariantTypeHandler) List(c *gin.Context) {
	types, err := h.svc.ListVariantTypes(c.Request.Context())

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	out := make([]VariantTypeResponse, len(types))
	for i, v := range types {
		out[i] = toVariantTypeResponse(v)
	}
	c.JSON(http.StatusOK, out)
}

// CreateType godoc
// @Summary     Create a variant type, e.g. Color (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body VariantTypeRequest true "Variant type"
// @Success     201 {object} VariantTypeResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/variant-types [post]
func (h *VariantTypeHandler) CreateType(c *gin.Context) {
	var req VariantTypeRequest
	if !httpx.Bind(c, &req) {
		return
	}

	t, err := h.svc.CreateVariantType(c.Request.Context(), usecase.VariantTypeInput{
		Name: req.Name, DataType: req.DataType, DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toVariantTypeResponse(*t))
}

// CreateValue godoc
// @Summary     Add a value to a variant type, e.g. Red (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Variant type ID"
// @Param       request body VariantValueRequest true "Variant value"
// @Success     201 {object} VariantValueResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/variant-types/{id}/values [post]
func (h *VariantTypeHandler) CreateValue(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	var req VariantValueRequest
	if !httpx.Bind(c, &req) {
		return
	}

	v, err := h.svc.CreateVariantValue(c.Request.Context(), id, usecase.VariantValueInput{
		Value: req.Value, DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}
	c.JSON(http.StatusCreated, VariantValueResponse{ID: v.ID, Value: v.Value, DisplayOrder: v.DisplayOrder})
}