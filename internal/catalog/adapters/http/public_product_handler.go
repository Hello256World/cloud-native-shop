package http

import (
	"context"
	"net/http"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/httpx"
	"github.com/gin-gonic/gin"
)

type publicProductService interface {
	ListPublicProducts(c context.Context, in usecase.ListProductsInput) (*usecase.ProductPage, error)
	GetPublicProductBySlug(c context.Context, slug string) (*domain.Product, error)
}

type PublicProductHandler struct {
	svc publicProductService
}

func NewPublicProductHandler(svc publicProductService) *PublicProductHandler {
	return &PublicProductHandler{
		svc: svc,
	}
}

// List godoc
// @Summary     List active products (cursor pagination)
// @Description Newest first. Pass the returned next_cursor as `cursor` to get the next page; null means last page.
// @Tags        catalog
// @Produce     json
// @Param       cursor      query string false "Opaque cursor from the previous page"
// @Param       limit       query int    false "Page size (1-50, default 20)"
// @Param       category_id query int    false "Filter by category id (exact, no descendants)"
// @Param       q           query string false "Search in product name"
// @Success     200 {object} ProductPageResponse
// @Failure     400 {object} httpx.ErrorBody
// @Router      /products [get]
func (h *PublicProductHandler) List(c *gin.Context) {
	var q ProductListQuery

	if !httpx.Bind(c, &q) {
		return
	}

	items, err := h.svc.ListPublicProducts(c.Request.Context(), toListProductsInput(q))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toProductPageResponse(items))
}

// Get godoc
// @Summary     Get an active product by slug (only active variants and images)
// @Tags        catalog
// @Produce     json
// @Param       slug path string true "Product slug"
// @Success     200 {object} ProductResponse
// @Failure     404 {object} httpx.ErrorBody
// @Router      /products/{slug} [get]
func (h *PublicProductHandler) Get(c *gin.Context) {
	items, err := h.svc.GetPublicProductBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		httpx.Respond(c, toAPIError(err))
		return
	}

	c.JSON(http.StatusOK, toProductResponse(items))
}
