package postgres

import (
	"context"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/ports"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type variantTypeRepository struct {
	pool *pgxpool.Pool
}

func NewVariantTypeRepository(pool *pgxpool.Pool) ports.VariantTypeRepository {
	return &variantTypeRepository{
		pool: pool,
	}
}

func (r *variantTypeRepository) q(c context.Context) *db.Queries {
	return db.New(postgres.Executor(c, r.pool))
}

func (r *variantTypeRepository) CreateType(c context.Context, v *domain.VariantType) (*domain.VariantType, error) {
	row, err := r.q(c).InsertVariantType(c, db.InsertVariantTypeParams{
		Name: v.Name, DataType: v.DataType,
		DisplayOrder: v.DisplayOrder,
	})

	if err != nil {
		if constraintContains(err, "name") {
			return nil, domain.ErrVariantTypeNameTaken
		}
		return nil, fmt.Errorf("catalog.repo.VariantType.CreateType: %w", err)
	}

	out := typeToDomain(row)
	return &out, nil
}

func (r *variantTypeRepository) GetTypeByID(c context.Context, id int64) (*domain.VariantType, error) {
	row, err := r.q(c).GetVariantTypeByID(c, id)

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrVariantTypeNotFound
		}
		return nil, fmt.Errorf("catalog.repo.VariantType.GetTypeByID id=%d: %w", id, err)
	}

	out := typeToDomain(row)
	return &out, nil
}

func (r *variantTypeRepository) ListTypes(c context.Context) ([]domain.VariantType, error) {
	q := r.q(c)

	types, err := q.ListVariantTypes(c)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.VariantType.ListTypes: %w", err)
	}

	values, err := q.ListVariantValues(c)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.VariantType.ListTypes values: %w", err)
	}

	out := make([]domain.VariantType, len(types))
	idx := make(map[int64]int, len(types))
	for i, val := range types {
		out[i] = typeToDomain(val)
		idx[val.ID] = i
	}

	for _, val := range values {
		if v, ok := idx[val.VariantTypeID]; ok {
			out[v].Values = append(out[v].Values, valueToDomain(val))
		}
	}

	return out, nil
}

func (r *variantTypeRepository) CreateValue(c context.Context, v *domain.VariantValue) (*domain.VariantValue, error) {
	row, err := r.q(c).InsertVariantValue(c, db.InsertVariantValueParams{
		VariantTypeID: v.TypeID,
		Value:         v.Value,
		DisplayOrder:  v.DisplayOrder,
	})

	if err != nil {
		if constraintContains(err, "type_id_value") {
			return nil, domain.ErrVariantValueTaken
		}
		return nil, fmt.Errorf("catalog.repo.VariantType.CreateValue: %w", err)
	}

	out := valueToDomain(row)
	return &out, nil
}

func (r *variantTypeRepository) ListValuesByIDs(c context.Context, ids []int64) ([]domain.VariantValue, error) {
	rows, err := r.q(c).ListVariantValuesByIDs(c, ids)

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.VariantType.ListValuesByIDs: %w", err)
	}

	out := make([]domain.VariantValue, len(rows))
	for i, v := range rows {
		out[i] = valueToDomain(v)
	}
	return out, nil
}
