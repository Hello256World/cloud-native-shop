package domain

import "errors"

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategorySlugTaken     = errors.New("category slug already exists")
	ErrCategoryHasChildren   = errors.New("category has children or products")
	ErrCategoryInvalidParent = errors.New("invalid parent category")

	ErrProductNotFound  = errors.New("product not found")
	ErrProductSlugTaken = errors.New("product slug already exists")

	ErrVariantNotFound = errors.New("variant not found")
	ErrSKUTaken        = errors.New("sku already exists")

	ErrProductInvalidCategory = errors.New("invalid product category")
	ErrInvalidVariantOptions  = errors.New("invalid variant options")
	ErrImageNotFound          = errors.New("image not found")

	ErrVariantTypeNotFound  = errors.New("variant type not found")
	ErrVariantTypeNameTaken = errors.New("variant type name already exists")
	ErrVariantValueTaken    = errors.New("variant value already exists")

	ErrInvalidCursor = errors.New("invalid cursor")
)
