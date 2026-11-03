package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type productSerivce interface {
	CreateProduct(context.Context, usecase.CreateProductInput) (*domain.Product, error)
	GetProduct(context.Context, int64) (*domain.Product, error)
	UpdateProduct(context.Context, int64, usecase.ProductInput) (*domain.Product, error)
	DeleteProduct(context.Context, int64) error
	AddVariant(ctx context.Context, productID int64, in usecase.VariantInput) (*domain.Variant, error)
	UpdateVariant(ctx context.Context, productID, variantID, price int64, active bool) (*domain.Variant, error)
	AddImage(ctx context.Context, productID int64, in usecase.ImageInput) (*domain.Image, error)
	DeleteImage(ctx context.Context, productID, imageID int64) error
}

type ProductHandler struct {
	svc productSerivce
}

func NewProductHandler(svc productSerivce) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

// Create godoc
// @Summary     Create a product with its variants and images in one transaction (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CreateProductRequest true "Product"
// @Success     201 {object} ProductResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     401 {object} httpx.ErrorBody
// @Failure     403 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if !httpx.Bind(c, &req) {
		return
	}

	p, err := h.svc.CreateProduct(c.Request.Context(), toCreateProductInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(p))
}

// Get godoc
// @Summary     Get a product by id, including inactive variants (admin)
// @Tags        admin-catalog
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Success     200 {object} ProductResponse
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	p, err := h.svc.GetProduct(c.Request.Context(), id)
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toProductResponse(p))
}

// Update godoc
// @Summary     Replace product fields; variants and images are untouched (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Param       request body ProductRequest true "Product"
// @Success     200 {object} ProductResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	var req ProductRequest
	if !httpx.Bind(c, &req) {
		return
	}

	p, err := h.svc.UpdateProduct(c.Request.Context(), id, toProductInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toProductResponse(p))
}

// Delete godoc
// @Summary     Soft-delete a product (admin)
// @Tags        admin-catalog
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Success     204
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	if err := h.svc.DeleteProduct(c.Request.Context(), id); err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// AddVariant godoc
// @Summary     Add a variant to a product (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Param       request body VariantRequest true "Variant"
// @Success     201 {object} VariantResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Failure     409 {object} httpx.ErrorBody
// @Router      /admin/products/{id}/variants [post]
func (h *ProductHandler) AddVariant(c *gin.Context) {
	id, ok := httpx.PathID(c)
	if !ok {
		return
	}

	var req VariantRequest
	if !httpx.Bind(c, &req) {
		return
	}

	v, err := h.svc.AddVariant(c.Request.Context(), id, toVariantInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toVariantResponse(*v))
}

// UpdateVariant godoc
// @Summary     Change price / active flag of a variant; SKU is immutable (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Param       variantId path int true "Variant ID"
// @Param       request body VariantUpdateRequest true "Variant"
// @Success     200 {object} VariantResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/products/{id}/variants/{variantId} [put]
func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	productID, ok := httpx.PathID(c)
	if !ok {
		return
	}

	variantID, ok := httpx.PathParamID(c, "variantId")
	if !ok {
		return
	}

	var req VariantUpdateRequest
	if !httpx.Bind(c, &req) {
		return
	}

	v, err := h.svc.UpdateVariant(c.Request.Context(), productID, variantID, req.Price, req.IsActive)
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toVariantResponse(*v))
}

// AddImage godoc
// @Summary     Add an image to a product (admin)
// @Tags        admin-catalog
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Param       request body ImageRequest true "Image"
// @Success     201 {object} ImageResponse
// @Failure     400 {object} httpx.ErrorBody
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/products/{id}/images [post]
func (h *ProductHandler) AddImage(c *gin.Context) {
	productID, ok := httpx.PathID(c)
	if !ok {
		return
	}

	var req ImageRequest
	if !httpx.Bind(c, &req) {
		return
	}

	i, err := h.svc.AddImage(c.Request.Context(), productID, toImageInput(req))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, toImageResponse(*i))
}

// DeleteImage godoc
// @Summary     Delete a product image (admin)
// @Tags        admin-catalog
// @Security    BearerAuth
// @Param       id path int true "Product ID"
// @Param       imageId path int true "Image ID"
// @Success     204
// @Failure     404 {object} httpx.ErrorBody
// @Router      /admin/products/{id}/images/{imageId} [delete]
func (h *ProductHandler) DeleteImage(c *gin.Context) {
	productID, ok := httpx.PathID(c)
	if !ok {
		return
	}
	imageID, ok := httpx.PathParamID(c, "imageId")
	if !ok {
		return
	}
	if err := h.svc.DeleteImage(c.Request.Context(), productID, imageID); err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}
	c.Status(http.StatusNoContent)
}
