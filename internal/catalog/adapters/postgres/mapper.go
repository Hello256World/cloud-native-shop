package postgres

import (
	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

func categoryToDomain(r db.CatalogCategory) *domain.Category {
	return &domain.Category{
		ID:        r.ID,
		ParentID:  r.ParentID,
		Name:      r.Name,
		Slug:      r.Slug,
		ImageName: r.ImageName,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func categoriesToDomain(rows []db.CatalogCategory) []domain.Category {
	out := make([]domain.Category, len(rows))
	for i, r := range rows {
		out[i] = *categoryToDomain(r)
	}
	return out
}

func productToDomain(p db.CatalogProduct) *domain.Product {
	return &domain.Product{
		ID:              p.ID,
		CategoryID:      p.CategoryID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Thumbnail:       p.Thumbnail,
		ShipmentWeightG: p.ShipmentWeightG,
		IsActive:        p.IsActive,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func variantToDomain(v db.CatalogProductVariant) domain.Variant {
	return domain.Variant{
		ID:        v.ID,
		ProductID: v.ProductID,
		SKU:       v.Sku,
		Price:     v.Price,
		IsActive:  v.IsActive,
	}
}

func imageToDomain(i db.CatalogProductImage) domain.Image {
	return domain.Image{
		ID:        i.ID,
		ImageName: i.ImageName,
		Priority:  i.Priority,
		IsActive:  i.IsActive,
	}
}

func typeToDomain(t db.CatalogVariantType) domain.VariantType {
	return domain.VariantType{ID: t.ID, Name: t.Name, DataType: t.DataType, DisplayOrder: t.DisplayOrder}
}

func valueToDomain(v db.CatalogVariantValue) domain.VariantValue {
	return domain.VariantValue{ID: v.ID, TypeID: v.VariantTypeID, Value: v.Value, DisplayOrder: v.DisplayOrder}
}
