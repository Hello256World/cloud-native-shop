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

func (r *attributeRepository) GetDefinitionByCode(c context.Context, code string) (*domain.AttributeDefinition, error) {
	row, err := r.q(c).GetAttributeDefinitionByCode(c, code)
	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrAttributeNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Attribute.GetDefinitionByCode code=%s: %w", code, err)
	}

	out := definitionToDomain(row)
	return &out, nil
}

func (r *attributeRepository) UpdateDefinition(c context.Context, d *domain.AttributeDefinition) (*domain.AttributeDefinition, error) {
	row, err := r.q(c).UpdateAttributeDefinition(c, db.UpdateAttributeDefinitionParams{
		ID:           d.ID,
		Name:         d.Name,
		Unit:         d.Unit,
		IsFilterable: d.IsFilterable,
		DisplayOrder: d.DisplayOrder,
	})
	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrAttributeNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Attribute.UpdateDefinition id=%d: %w", d.ID, err)
	}
	out := definitionToDomain(row)
	return &out, nil
}

func (r *attributeRepository) ListFacets(ctx context.Context, categoryID int64) ([]domain.AttributeFacet, error) {
	rows, err := r.q(ctx).ListAttributeFacets(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Attribute.ListFacets: %w", err)
	}

	out := make([]domain.AttributeFacet, 0)
	idx := make(map[int64]int)
	for _, v := range rows {
		i, ok := idx[v.AttributeID]
		if !ok {
			out = append(out, domain.AttributeFacet{Attribute: domain.AttributeDefinition{
				ID: v.AttributeID, Code: v.Code, Name: v.Name,
				ValueType: domain.AttributeValueType(v.ValueType),
				Unit:      v.Unit, IsFilterable: true, DisplayOrder: v.DisplayOrder,
			}})
			i = len(out) - 1
			idx[v.AttributeID] = i
		}
		out[i].Values = append(out[i].Values, domain.FacetValue{
			ValueText: v.ValueText, ValueNumber: v.ValueNumber, Count: v.ProductCount,
		})
	}
	return out, nil
}
