package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type categoryService interface {
	CreateCategory(context.Context, usecase.CategoryInput) (*domain.Category, error)
	UpdateCategory(context.Context, int64, usecase.CategoryInput) (*domain.Category, error)
	DeleteCategory(context.Context, int64) error
	ListCategories(context.Context) ([]domain.Category, error)
	ListPublicCategories(context.Context) ([]domain.Category, error)
}

type CategoryHandler struct {
	svc categoryService
}

func NewCategoryHandler(svc categoryService) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
	}
}

// ListPublic godoc
// @Summary     List active categories (flat; build the tree with parent_id)
// @Tags        catalog
// @Produce     json
// @Success     200 {array} CategoryResponse
// @Router      /categories [get]
func (h *CategoryHandler) ListPublic(c *gin.Context) {
	list, err := h.svc.ListPublicCategories(c.Request.Context())

	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toCategoryListResponse(list))
}

// ListAdmin godoc
// @Summary     List all categories (admin)
// @Tags        admin-catalog
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} CategoryResponse
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Router      /admin/categories [get]
func (h *CategoryHandler) ListAdmin(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toCategoryListResponse(list))
}

// Create godoc
// @Summary     Create a category (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CategoryRequest true "Category"
// @Success     201 {object} CategoryResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CategoryRequest
	if !httpx.Bind(c, &req) {
		return
	}

	cat, err := h.svc.CreateCategory(c.Request.Context(), toCategoryInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toCategoryResponse(cat))
}

// Update godoc
// @Summary     Replace a category (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Category ID"
// @Param       request body CategoryRequest true "Category"
// @Success     200 {object} CategoryResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}
	var req CategoryRequest
	if !httpx.Bind(c, &req) {
		return
	}

	cat, err := h.svc.UpdateCategory(c.Request.Context(), id, toCategoryInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toCategoryResponse(cat))
}

// Delete godoc
// @Summary     Delete a category (admin, must have no subcategories or products)
// @Tags        admin-catalog
// @Security    BearerAuth
// @Param       id path int true "Category ID"
// @Success     204
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}
	c.Status(http.StatusNoContent)
}
