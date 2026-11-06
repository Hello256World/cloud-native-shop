package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/ports"
)

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r *productRepository) ListPublic(c context.Context, f ports.ProductListFilter) ([]domain.ProductSummary, error) {
	rows, err := r.q(c).ListPublicProducts(c, db.ListPublicProductsParams{
		AfterID:    f.AfterID,
		CategoryID: f.CategoryID,
		Search:     likeEscaper.Replace(f.Search),
		PageSize:   f.Limit,
	})

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.ListPublic: %w", err)
	}

	out := make([]domain.ProductSummary, len(rows))
	for i, v := range rows {
		out[i] = domain.ProductSummary{
			ID:         v.ID,
			CategoryID: v.CategoryID,
			Name:       v.Name,
			Slug:       v.Slug,
			Thumbnail:  v.Thumbnail,
			MinPrice:   v.MinPrice,
			IsActive:   v.IsActive,
		}
	}

	return out, nil
}

func (r *productRepository) ListAdmin(c context.Context, f ports.ProductListFilter) ([]domain.ProductSummary, error) {
	rows, err := r.q(c).ListAdminProducts(c, db.ListAdminProductsParams{
		AfterID:    f.AfterID,
		CategoryID: f.CategoryID,
		Search:     likeEscaper.Replace(f.Search),
		PageSize:   f.Limit,
	})

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.ListAdmin: %w", err)
	}

	out := make([]domain.ProductSummary, len(rows))
	for i, v := range rows {
		out[i] = domain.ProductSummary{
			ID: v.ID, CategoryID: v.CategoryID, Name: v.Name, Slug: v.Slug,
			Thumbnail: v.Thumbnail, MinPrice: v.MinPrice, IsActive: v.IsActive,
		}
	}
	return out, nil
}
