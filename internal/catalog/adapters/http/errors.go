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
	default:
		return err
	}
	return apierrors.Wrap(api, err)
}
