package usecase

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type VariantTypeInput struct {
	Name         string
	DataType     string
	DisplayOrder int
}

type VariantValueInput struct {
	Value        string
	DisplayOrder int
}

func (s *Service) CreateVariantType(c context.Context, in VariantTypeInput) (*domain.VariantType, error) {
	return s.types.CreateType(c, &domain.VariantType{
		Name: in.Name, DataType: in.DataType, DisplayOrder: in.DisplayOrder,
	})
}

func (s *Service) ListVariantTypes(c context.Context) ([]domain.VariantType, error) {
	return s.types.ListTypes(c)
}

func (s *Service) CreateVariantValue(c context.Context, typeID int64, in VariantValueInput) (*domain.VariantValue, error) {
	_, err := s.types.GetTypeByID(c, typeID)
	if err != nil {
		return nil, err
	}

	return s.types.CreateValue(c, &domain.VariantValue{
		TypeID:       typeID,
		Value:        in.Value,
		DisplayOrder: in.DisplayOrder,
	})
}