package usecase

import (
	"context"
	"errors"
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

type AttributeDefinitionUpdateInput struct {
	Name         string
	Unit         *string
	IsFilterable bool
	DisplayOrder int
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

func (s *Service) UpdateAttributeDefinition(c context.Context, id int64, in AttributeDefinitionUpdateInput) (*domain.AttributeDefinition, error) {
	return s.attrs.UpdateDefinition(c, &domain.AttributeDefinition{
		ID: id, Name: in.Name, Unit: in.Unit,
		IsFilterable: in.IsFilterable, DisplayOrder: in.DisplayOrder,
	})
}

func (s *Service) ListProductFacets(c context.Context, categoryID *int64) ([]domain.AttributeFacet, error) {
	var id int64
	if categoryID != nil {
		id = *categoryID
	}
	return s.attrs.ListFacets(c, id)
}

func (s *Service) resolveAttrFilter(c context.Context, code, value string) (id int64, num *float64, txt *string, err error) {
	def, err := s.attrs.GetDefinitionByCode(c, strings.TrimSpace(code))
	if err != nil {
		if errors.Is(err, domain.ErrAttributeNotFound) {
			return 0, nil, nil, domain.ErrInvalidAttrFilter
		}
		return 0, nil, nil, err
	}
	if !def.IsFilterable {
		return 0, nil, nil, domain.ErrInvalidAttrFilter
	}
	num, txt, perr := parseAttrValue(*def, value)
	if perr != nil {
		return 0, nil, nil, domain.ErrInvalidAttrFilter
	}
	return def.ID, num, txt, nil
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
		num, txt, err := parseAttrValue(def, a.Value)
		if err != nil {
			return nil, err
		}
		out = append(out, domain.ProductAttribute{AttributeID: def.ID, ValueNumber: num, ValueText: txt})
	}
	return out, nil
}

func parseAttrValue(def domain.AttributeDefinition, raw string) (num *float64, txt *string, err error) {
	val := strings.TrimSpace(raw)
	if val == "" {
		return nil, nil, domain.ErrInvalidAttributes
	}
	if def.ValueType == domain.AttrNumber {
		n, perr := strconv.ParseFloat(val, 64)
		if perr != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, nil, domain.ErrInvalidAttributes
		}
		return &n, nil, nil
	}
	return nil, &val, nil
}
