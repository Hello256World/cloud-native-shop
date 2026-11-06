package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type ProductListFilter struct {
	AfterID    int64
	CategoryID int64
	Search     string
	Limit      int
}

type ProductRepository interface {
	Create(context.Context, *domain.Product) (*domain.Product, error)
	GetByID(context.Context, int64) (*domain.Product, error)
	Update(context.Context, *domain.Product) (*domain.Product, error)
	Delete(context.Context, int64) error

	AddVariant(ctx context.Context, productID int64, v *domain.Variant) (*domain.Variant, error)
	UpdateVariant(ctx context.Context, v *domain.Variant) (*domain.Variant, error)

	AddImage(ctx context.Context, productID int64, img *domain.Image) (*domain.Image, error)
	DeleteImage(ctx context.Context, productID, imageID int64) error

	GetBySlug(context.Context, string) (*domain.Product, error)
	ListPublic(context.Context, ProductListFilter) ([]domain.ProductSummary, error)
	ListAdmin(context.Context, ProductListFilter) ([]domain.ProductSummary, error)
}
