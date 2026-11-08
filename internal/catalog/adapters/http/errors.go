package http

import (
	"errors"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
)

func toAPIError(err error) error {
	var api apierrors.APIError

	switch {
	case errors.Is(err, domain.ErrCategoryNotFound):
		api = apierrors.NewNotFoundError("category_not_found", "category not found")
	case errors.Is(err, domain.ErrCategorySlugTaken):
		api = apierrors.NewConflictError("category_slug_taken", "category slug already exists")
	case errors.Is(err, domain.ErrCategoryHasChildren):
		api = apierrors.NewConflictError("category_not_empty", "category has subcategories or products")
	case errors.Is(err, domain.ErrCategoryInvalidParent):
		api = apierrors.NewBadRequestError("invalid_parent", "invalid parent category")
	case errors.Is(err, domain.ErrProductNotFound):
		api = apierrors.NewNotFoundError("product_not_found", "product not found")
	case errors.Is(err, domain.ErrProductSlugTaken):
		api = apierrors.NewConflictError("product_slug_taken", "product slug already exists")
	case errors.Is(err, domain.ErrProductInvalidCategory):
		api = apierrors.NewBadRequestError("invalid_category", "invalid category")
	case errors.Is(err, domain.ErrVariantNotFound):
		api = apierrors.NewNotFoundError("variant_not_found", "variant not found")
	case errors.Is(err, domain.ErrSKUTaken):
		api = apierrors.NewConflictError("sku_taken", "sku already exists")
	case errors.Is(err, domain.ErrInvalidVariantOptions):
		api = apierrors.NewBadRequestError("invalid_variant_options", "invalid variant options")
	case errors.Is(err, domain.ErrImageNotFound):
		api = apierrors.NewNotFoundError("image_not_found", "image not found")
	case errors.Is(err, domain.ErrVariantTypeNotFound):
		api = apierrors.NewNotFoundError("variant_type_not_found", "variant type not found")
	case errors.Is(err, domain.ErrVariantTypeNameTaken):
		api = apierrors.NewConflictError("variant_type_name_taken", "variant type name already exists")
	case errors.Is(err, domain.ErrVariantValueTaken):
		api = apierrors.NewConflictError("variant_value_taken", "variant value already exists")
	case errors.Is(err, domain.ErrInvalidCursor):
		api = apierrors.NewBadRequestError("invalid_cursor", "invalid cursor")
	case errors.Is(err, domain.ErrInvalidSort):
		api = apierrors.NewBadRequestError("invalid_sort", "invalid sort")
	default:
		return err
	}
	return apierrors.Wrap(api, err)
}
