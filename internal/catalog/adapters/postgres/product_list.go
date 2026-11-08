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

func summary(id, categoryID int64, name, slug, thumb string, price int64, active bool) domain.ProductSummary {
	return domain.ProductSummary{
		ID: id, CategoryID: categoryID, Name: name, Slug: slug,
		Thumbnail: thumb, MinPrice: price, IsActive: active,
	}
}

func (r *productRepository) ListPublic(c context.Context, f ports.ProductListFilter) ([]domain.ProductSummary, error) {
	q := r.q(c)
	search := likeEscaper.Replace(f.Search)

	var out []domain.ProductSummary

	switch f.Sort {
	case domain.SortPriceAsc:
		rows, err := q.ListPublicProductsByPriceAsc(c, db.ListPublicProductsByPriceAscParams{
			AfterPrice: f.AfterPrice, AfterID: f.AfterID, CategoryID: f.CategoryID,
			Search: search, PriceMin: f.PriceMin, PriceMax: f.PriceMax, PageSize: f.Limit,
		})
		if err != nil {
			return nil, fmt.Errorf("catalog.repo.Product.ListPublic price_asc: %w", err)
		}
		for _, v := range rows {
			out = append(out, summary(v.ID, v.CategoryID, v.Name, v.Slug, v.Thumbnail, v.MinPrice, v.IsActive))
		}

	case domain.SortPriceDesc:
		rows, err := q.ListPublicProductsByPriceDesc(c, db.ListPublicProductsByPriceDescParams{
			AfterPrice: f.AfterPrice, AfterID: f.AfterID, CategoryID: f.CategoryID,
			Search: search, PriceMin: f.PriceMin, PriceMax: f.PriceMax, PageSize: f.Limit,
		})
		if err != nil {
			return nil, fmt.Errorf("catalog.repo.Product.ListPublic price_desc: %w", err)
		}
		for _, v := range rows {
			out = append(out, summary(v.ID, v.CategoryID, v.Name, v.Slug, v.Thumbnail, v.MinPrice, v.IsActive))
		}

	default: 
		rows, err := q.ListPublicProducts(c, db.ListPublicProductsParams{
			AfterID: f.AfterID, CategoryID: f.CategoryID, Search: search,
			PriceMin: f.PriceMin, PriceMax: f.PriceMax, PageSize: f.Limit,
		})
		if err != nil {
			return nil, fmt.Errorf("catalog.repo.Product.ListPublic newest: %w", err)
		}
		for _, v := range rows {
			out = append(out, summary(v.ID, v.CategoryID, v.Name, v.Slug, v.Thumbnail, v.MinPrice, v.IsActive))
		}
	}

	if out == nil {
		out = []domain.ProductSummary{}
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
