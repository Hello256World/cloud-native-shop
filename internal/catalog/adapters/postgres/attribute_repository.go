package postgres

import (
	"context"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/ports"
	pgTX "github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type attributeRepository struct {
	pool *pgxpool.Pool
}

func NewAttributeRepository(pool *pgxpool.Pool) ports.AttributeRepository {
	return &attributeRepository{
		pool: pool,
	}
}

func (r *attributeRepository) q(c context.Context) *db.Queries {
	return db.New(pgTX.Executor(c, r.pool))
}

func (r *attributeRepository) CreateDefinition(c context.Context, def *domain.AttributeDefinition) (*domain.AttributeDefinition, error) {
	row, err := r.q(c).InsertAttributeDefinition(c, db.InsertAttributeDefinitionParams{
		Code:         def.Code,
		Name:         def.Name,
		ValueType:    string(def.ValueType),
		Unit:         def.Unit,
		IsFilterable: def.IsFilterable,
		DisplayOrder: def.DisplayOrder,
	})

	if err != nil {
		if constraintContains(err, "code") {
			return nil, domain.ErrAttributeCodeTaken
		}
		return nil, fmt.Errorf("catalog.repo.Attribute.CreateDefinition: %w", err)
	}

	out := definitionToDomain(row)
	return &out, nil
}

func (r *attributeRepository) ListDefinitions(c context.Context) ([]domain.AttributeDefinition, error) {
	rows, err := r.q(c).ListAttributeDefinitions(c)

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Attribute.ListDefinitions: %w", err)
	}

	out := make([]domain.AttributeDefinition, len(rows))
	for i, v := range rows {
		out[i] = definitionToDomain(v)
	}

	return out, nil
}

func (r *attributeRepository) ListDefinitionsByIDs(c context.Context, ids []int64) ([]domain.AttributeDefinition, error) {
	rows, err := r.q(c).ListAttributeDefinitionsByIDs(c, ids)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Attribute.ListDefinitionsByIDs: %w", err)
	}

	out := make([]domain.AttributeDefinition, len(rows))
	for i, v := range rows {
		out[i] = definitionToDomain(v)
	}

	return out, nil
}
