package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type ProductRepository interface {
	Create(context.Context, *domain.Product) (*domain.Product, error)
	GetByID(context.Context, int64) (*domain.Product, error)
	Update(context.Context, *domain.Product) (*domain.Product, error)
	Delete(context.Context, int64) error

	AddVariant(ctx context.Context, productID int64, v *domain.Variant) (*domain.Variant, error)
	UpdateVariant(ctx context.Context, v *domain.Variant) (*domain.Variant, error)

	AddImage(ctx context.Context, productID int64, img *domain.Image) (*domain.Image, error)
	DeleteImage(ctx context.Context, productID, imageID int64) error
}
