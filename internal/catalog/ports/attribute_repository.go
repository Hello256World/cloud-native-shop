package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type AttributeRepository interface {
	CreateDefinition(context.Context, *domain.AttributeDefinition) (*domain.AttributeDefinition, error)
	ListDefinitions(context.Context) ([]domain.AttributeDefinition, error)
	ListDefinitionsByIDs(context.Context, []int64) ([]domain.AttributeDefinition, error)
	GetDefinitionByCode(context.Context, string) (*domain.AttributeDefinition, error)
	UpdateDefinition(context.Context, *domain.AttributeDefinition) (*domain.AttributeDefinition, error)
	ListFacets(ctx context.Context, categoryID int64) ([]domain.AttributeFacet, error)
}
