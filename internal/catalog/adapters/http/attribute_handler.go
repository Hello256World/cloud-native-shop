package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type attributeService interface {
	CreateAttributeDefinition(context.Context, usecase.AttributeDefinitionInput) (*domain.AttributeDefinition, error)
	ListAttributeDefinitions(context.Context) ([]domain.AttributeDefinition, error)
	UpdateAttributeDefinition(context.Context, int64, usecase.AttributeDefinitionUpdateInput) (*domain.AttributeDefinition, error)
}

type AttributeHandler struct {
	svc attributeService
}

func NewAttributeHandler(svc attributeService) *AttributeHandler {
	return &AttributeHandler{
		svc: svc,
	}
}

// Create godoc
// @Summary     Create an attribute definition, e.g. RAM (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body AttributeDefinitionRequest true "Attribute definition"
// @Success     201 {object} AttributeDefinitionResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/attributes [post]
func (h *AttributeHandler) Create(c *gin.Context) {
	var req AttributeDefinitionRequest
	if !httpx.Bind(c, &req) {
		return
	}

	a, err := h.svc.CreateAttributeDefinition(c.Request.Context(), usecase.AttributeDefinitionInput{
		Code:         req.Code,
		Name:         req.Name,
		ValueType:    domain.AttributeValueType(req.ValueType),
		Unit:         req.Unit,
		IsFilterable: req.IsFilterable,
		DisplayOrder: req.DisplayOrder,
	})

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toAttributeDefinitionResponse(*a))
}

// List godoc
// @Summary     List attribute definitions (admin)
// @Tags        admin-catalog
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} AttributeDefinitionResponse
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Router      /admin/attributes [get]
func (h *AttributeHandler) List(c *gin.Context) {
	defs, err := h.svc.ListAttributeDefinitions(c.Request.Context())
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	out := make([]AttributeDefinitionResponse, len(defs))
	for i, v := range defs {
		out[i] = toAttributeDefinitionResponse(v)
	}
	c.JSON(http.StatusOK, out)
}

// Update godoc
// @Summary     Update an attribute definition; code and value_type are immutable (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Attribute ID"
// @Param       request body AttributeDefinitionUpdateRequest true "Attribute definition"
// @Success     200 {object} AttributeDefinitionResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/attributes/{id} [put]
func (h *AttributeHandler) Update(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	var req AttributeDefinitionUpdateRequest
	if !httpx.Bind(c, &req) {
		return
	}

	d, err := h.svc.UpdateAttributeDefinition(c.Request.Context(), id, usecase.AttributeDefinitionUpdateInput{
		Name: req.Name, Unit: req.Unit, IsFilterable: req.IsFilterable, DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toAttributeDefinitionResponse(*d))
}
