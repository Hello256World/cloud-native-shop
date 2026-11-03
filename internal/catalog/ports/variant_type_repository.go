package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type VariantTypeRepository interface {
	CreateType(context.Context, *domain.VariantType) (*domain.VariantType, error)
	GetTypeByID(context.Context, int64) (*domain.VariantType, error)
	ListTypes(context.Context) ([]domain.VariantType, error)
	CreateValue(context.Context, *domain.VariantValue) (*domain.VariantValue, error)
	ListValuesByIDs(context.Context, []int64) ([]domain.VariantValue, error)
}
