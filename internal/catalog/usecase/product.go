package usecase

import (
	"context"
	"errors"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type ProductInput struct {
	CategoryID      int64
	Name            string
	Slug            string
	Description     *string
	Thumbnail       string
	ShipmentWeightG int
	IsActive        bool
}

type VariantInput struct {
	SKU      string
	Price    int64
	IsActive bool
	ValueIDs []int64
}

type ImageInput struct {
	ImageName string
	Priority  int
	IsActive  bool
}

type CreateProductInput struct {
	ProductInput
	Variants []VariantInput
	Images   []ImageInput
}

func (in ProductInput) toDomain(id int64) *domain.Product {
	return &domain.Product{
		ID:              id,
		CategoryID:      in.CategoryID,
		Name:            in.Name,
		Slug:            in.Slug,
		Description:     in.Description,
		Thumbnail:       in.Thumbnail,
		ShipmentWeightG: in.ShipmentWeightG,
		IsActive:        in.IsActive,
	}
}

func (in ImageInput) toDomain() *domain.Image {
	return &domain.Image{
		ImageName: in.ImageName, Priority: in.Priority, IsActive: in.IsActive,
	}
}

func (in VariantInput) toDomain(opts []domain.VariantOption) *domain.Variant {
	return &domain.Variant{
		SKU:      in.SKU,
		Price:    in.Price,
		IsActive: in.IsActive,
		Options:  opts,
	}
}

func (s *Service) CreateProduct(c context.Context, in CreateProductInput) (*domain.Product, error) {
	if err := s.checkCategoryID(c, in.CategoryID); err != nil {
		return nil, err
	}

	opts := make([][]domain.VariantOption, len(in.Variants))
	for i, val := range in.Variants {
		opt, err := s.buildOption(c, val.ValueIDs)
		if err != nil {
			return nil, err
		}
		opts[i] = opt
	}

	var out *domain.Product
	if err := s.tx.WithinTx(c, func(ctx context.Context) error {
		p, err := s.products.Create(ctx, in.toDomain(0))
		if err != nil {
			return err
		}

		for i, val := range in.Variants {
			if _, err := s.products.AddVariant(ctx, p.ID, val.toDomain(opts[i])); err != nil {
				return err
			}
		}

		for _, val := range in.Images {
			if _, err := s.products.AddImage(ctx, p.ID, val.toDomain()); err != nil {
				return err
			}
		}

		out, err = s.products.GetByID(ctx, p.ID)
		return err
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetProduct(c context.Context, id int64) (*domain.Product, error) {
	return s.products.GetByID(c, id)
}

func (s *Service) UpdateProduct(c context.Context, id int64, in ProductInput) (*domain.Product, error) {
	if _, err := s.products.GetByID(c, id); err != nil {
		return nil, err
	}

	if err := s.checkCategoryID(c, in.CategoryID); err != nil {
		return nil, err
	}

	if _, err := s.products.Update(c, in.toDomain(id)); err != nil {
		return nil, err
	}

	return s.products.GetByID(c, id)
}

func (s *Service) DeleteProduct(c context.Context, id int64) error {
	return s.products.Delete(c, id)
}

func (s *Service) AddVariant(c context.Context, productID int64, in VariantInput) (*domain.Variant, error) {
	if _, err := s.products.GetByID(c, productID); err != nil {
		return nil, err
	}

	opts, err := s.buildOption(c, in.ValueIDs)
	if err != nil {
		return nil, err
	}

	return s.products.AddVariant(c, productID, in.toDomain(opts))
}

func (s *Service) UpdateVariant(c context.Context, productID, variantID, price int64, actice bool) (*domain.Variant, error) {
	return s.products.UpdateVariant(c, &domain.Variant{
		ProductID: productID,
		ID:        variantID,
		Price:     price,
		IsActive:  actice,
	})
}

func (s *Service) AddImage(c context.Context, productID int64, in ImageInput) (*domain.Image, error) {
	if _, err := s.products.GetByID(c, productID); err != nil {
		return nil, err
	}

	return s.products.AddImage(c, productID, in.toDomain())
}

func (s *Service) DeleteImage(c context.Context, productID, imageID int64) error {
	return s.products.DeleteImage(c, productID, imageID)
}

func (s *Service) checkCategoryID(c context.Context, id int64) error {
	if _, err := s.categories.GetByID(c, id); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return domain.ErrProductInvalidCategory
		}
		return err
	}

	return nil
}

func (s *Service) buildOption(c context.Context, ids []int64) ([]domain.VariantOption, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	values, err := s.types.ListValuesByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	if len(values) != len(ids) {
		return nil, domain.ErrInvalidVariantOptions
	}

	seenType := make(map[int64]bool, len(values))
	opts := make([]domain.VariantOption, 0, len(values))
	for _, v := range values {
		if seenType[v.TypeID] {
			return nil, domain.ErrInvalidVariantOptions
		}
		seenType[v.TypeID] = true
		opts = append(opts, domain.VariantOption{TypeID: v.TypeID, ValueID: v.ID})
	}
	return opts, nil
}
