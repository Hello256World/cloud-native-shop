package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type CategoryRepository interface {
	Create(context.Context, *domain.Category) (*domain.Category, error)
	GetByID(context.Context, int64) (*domain.Category, error)
	ListAll(context.Context) ([]domain.Category, error)
	ListActive(context.Context) ([]domain.Category, error)
	Update(context.Context, *domain.Category) (*domain.Category, error)
	Delete(context.Context, int64) error
	HasChildren(context.Context, int64) (bool, error)
	HasProducts(context.Context, int64) (bool, error)
}
