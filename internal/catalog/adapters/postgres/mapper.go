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
