package usecase

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type AttributeDefinitionInput struct {
	Code         string
	Name         string
	ValueType    domain.AttributeValueType
	Unit         *string
	IsFilterable bool
	DisplayOrder int
}

type AttributeInput struct {
	AttributeID int64
	Value       string
}

func (s *Service) CreateAttributeDefinition(c context.Context, in AttributeDefinitionInput) (*domain.AttributeDefinition, error) {
	if !in.ValueType.IsValid() {
		return nil, domain.ErrInvalidAttributes
	}
	return s.attrs.CreateDefinition(c, &domain.AttributeDefinition{
		Code:         in.Code,
		Name:         in.Name,
		ValueType:    in.ValueType,
		Unit:         in.Unit,
		IsFilterable: in.IsFilterable,
		DisplayOrder: in.DisplayOrder,
	})
}

func (s *Service) ListAttributeDefinitions(c context.Context) ([]domain.AttributeDefinition, error) {
	return s.attrs.ListDefinitions(c)
}

func (s *Service) buildAttributes(c context.Context, in []AttributeInput) ([]domain.ProductAttribute, error) {
	if len(in) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(in))
	seen := make(map[int64]bool, len(in))
	for _, a := range in {
		if seen[a.AttributeID] {
			return nil, domain.ErrInvalidAttributes
		}
		seen[a.AttributeID] = true
		ids = append(ids, a.AttributeID)
	}

	defs, err := s.attrs.ListDefinitionsByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	byID := make(map[int64]domain.AttributeDefinition, len(defs))
	for _, d := range defs {
		byID[d.ID] = d
	}

	out := make([]domain.ProductAttribute, 0, len(in))
	for _, a := range in {
		def, ok := byID[a.AttributeID]
		if !ok {
			return nil, domain.ErrInvalidAttributes
		}

		val := strings.TrimSpace(a.Value)
		if val == "" {
			return nil, domain.ErrInvalidAttributes
		}

		pa := domain.ProductAttribute{AttributeID: def.ID}
		if def.ValueType == domain.AttrNumber {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
				return nil, domain.ErrInvalidAttributes
			}
			pa.ValueNumber = &n
		} else {
			pa.ValueText = &val
		}
		out = append(out, pa)
	}
	return out, nil
}
