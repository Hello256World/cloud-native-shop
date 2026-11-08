package usecase

import (
	"context"
	"strings"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/ports"
)

const (
	defaultPageSize = 20
	maxPageSize     = 50
)

type ListProductsInput struct {
	Cursor     string
	Limit      int
	CategoryID *int64
	Search     string
	Sort       domain.ProductSort
	PriceMin   int64
	PriceMax   int64
}

type ProductPage struct {
	Items      []domain.ProductSummary
	NextCursor string
}

type listFunc func(context.Context, ports.ProductListFilter) ([]domain.ProductSummary, error)

func (s *Service) ListPublicProducts(c context.Context, in ListProductsInput) (*ProductPage, error) {
	return s.listProducts(c, in, s.products.ListPublic)
}

func (s *Service) ListAdminProducts(c context.Context, in ListProductsInput) (*ProductPage, error) {
	in.Sort = domain.SortNewest
	in.PriceMin, in.PriceMax = 0, 0
	return s.listProducts(c, in, s.products.ListAdmin)
}

func (s *Service) listProducts(c context.Context, in ListProductsInput, fetch listFunc) (*ProductPage, error) {
	sort := in.Sort
	if sort == "" {
		sort = domain.SortNewest
	}
	if !sort.Valid() {
		return nil, domain.ErrInvalidSort
	}

	afterPrice, afterID, err := decodeCursor(in.Cursor, sort)
	if err != nil {
		return nil, err
	}

	limit := in.Limit
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	f := ports.ProductListFilter{
		Sort:       sort,
		AfterID:    afterID,
		AfterPrice: afterPrice,
		Search:     strings.TrimSpace(in.Search),
		PriceMin:   in.PriceMin,
		PriceMax:   in.PriceMax,
		Limit:      limit + 1,
	}
	if in.CategoryID != nil {
		f.CategoryID = *in.CategoryID
	}

	items, err := fetch(c, f)
	if err != nil {
		return nil, err
	}

	page := &ProductPage{Items: items}
	if len(items) > limit {
		page.Items = items[:limit]
		last := items[limit-1]
		var price int64
		if sort != domain.SortNewest {
			price = last.MinPrice
		}
		page.NextCursor = encodeCursor(sort, price, last.ID)
	}
	return page, nil
}

func (s *Service) GetPublicProductBySlug(c context.Context, slug string) (*domain.Product, error) {
	p, err := s.products.GetBySlug(c, slug)
	if err != nil {
		return nil, err
	}
	if !p.IsActive {
		return nil, domain.ErrProductNotFound
	}

	variants := make([]domain.Variant, 0, len(p.Variants))
	for _, v := range p.Variants {
		if v.IsActive {
			variants = append(variants, v)
		}
	}
	if len(variants) == 0 {
		return nil, domain.ErrProductNotFound
	}
	p.Variants = variants

	images := make([]domain.Image, 0, len(p.Images))
	for _, im := range p.Images {
		if im.IsActive {
			images = append(images, im)
		}
	}
	p.Images = images

	return p, nil
}
