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

type productRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) ports.ProductRepository {
	return &productRepository{
		pool: pool,
	}
}

func (r *productRepository) q(c context.Context) *db.Queries {
	return db.New(postgres.Executor(c, r.pool))
}

func (r *productRepository) Create(c context.Context, p *domain.Product) (*domain.Product, error) {
	row, err := r.q(c).InsertProduct(c, db.InsertProductParams{
		CategoryID:      p.CategoryID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Thumbnail:       p.Thumbnail,
		ShipmentWeightG: p.ShipmentWeightG,
		IsActive:        p.IsActive,
	})

	if err != nil {
		if constraintContains(err, "slug") {
			return nil, domain.ErrProductSlugTaken
		}
		return nil, fmt.Errorf("catalog.repo.Product.Create: %w", err)
	}
	return productToDomain(row), nil
}

func (r *productRepository) GetByID(c context.Context, id int64) (*domain.Product, error) {
	q := r.q(c)

	row, err := q.GetProductByID(c, id)

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Product.GetByID id=%d: %w", id, err)
	}

	return r.assemble(c, q, row)
}

func (r *productRepository) Update(c context.Context, product *domain.Product) (*domain.Product, error) {
	q := r.q(c)

	row, err := q.UpdateProduct(c, db.UpdateProductParams{
		ID:              product.ID,
		CategoryID:      product.CategoryID,
		Name:            product.Name,
		Slug:            product.Slug,
		Description:     product.Description,
		Thumbnail:       product.Thumbnail,
		ShipmentWeightG: product.ShipmentWeightG,
		IsActive:        product.IsActive,
	})

	if err != nil {
		switch {
		case isNoRows(err):
			return nil, domain.ErrProductNotFound
		case constraintContains(err, "slug"):
			return nil, domain.ErrProductSlugTaken
		}
		return nil, fmt.Errorf("catalog.repo.Product.Update id=%d: %w", product.ID, err)
	}

	return productToDomain(row), nil
}

func (r *productRepository) Delete(c context.Context, id int64) error {
	n, err := r.q(c).SoftDeleteProduct(c, id)
	if err != nil {
		return fmt.Errorf("catalog.repo.Product.Delete id=%d: %w", id, err)
	}
	if n == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) AddVariant(ctx context.Context, productID int64, v *domain.Variant) (*domain.Variant, error) {
	q := r.q(ctx)

	row, err := q.InsertVariant(ctx, db.InsertVariantParams{
		ProductID: productID,
		Sku:       v.SKU,
		Price:     v.Price,
		IsActive:  v.IsActive,
	})

	if err != nil {
		if constraintContains(err, "sku") {
			return nil, domain.ErrSKUTaken
		}
		return nil, fmt.Errorf("catalog.repo.Product.AddVariant product=%d: %w", productID, err)
	}

	for _, o := range v.Options {
		if err := q.InsertVariantOption(ctx, db.InsertVariantOptionParams{
			ProductVariantID: row.ID, VariantValueID: o.ValueID,
		}); err != nil {
			return nil, fmt.Errorf("catalog.repo.Product.AddVariant option variant=%d: %w", row.ID, err)
		}
	}

	out := variantToDomain(row)
	out.Options = v.Options
	return &out, nil
}

func (r *productRepository) UpdateVariant(ctx context.Context, v *domain.Variant) (*domain.Variant, error) {
	q := r.q(ctx)

	row, err := q.UpdateVariant(ctx, db.UpdateVariantParams{
		ID:        v.ID,
		ProductID: v.ProductID,
		Price:     v.Price,
		IsActive:  v.IsActive,
	})

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrVariantNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Product.UpdateVariant id=%d: %w", v.ID, err)
	}

	out := variantToDomain(row)
	return &out, nil
}

func (r *productRepository) AddImage(ctx context.Context, productID int64, img *domain.Image) (*domain.Image, error) {
	q := r.q(ctx)

	row, err := q.InsertImage(ctx, db.InsertImageParams{
		ProductID: productID,
		ImageName: img.ImageName,
		Priority:  img.Priority,
		IsActive:  img.IsActive,
	})

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.AddImage product=%d: %w", productID, err)
	}

	out := imageToDomain(row)
	return &out, nil
}

func (r *productRepository) DeleteImage(ctx context.Context, productID, imageID int64) error {
	n, err := r.q(ctx).DeleteImage(ctx, db.DeleteImageParams{
		ID:        imageID,
		ProductID: productID,
	})
	if err != nil {
		return fmt.Errorf("catalog.repo.Product.DeleteImage id=%d: %w", imageID, err)
	}
	if n == 0 {
		return domain.ErrImageNotFound
	}
	return nil
}

func (r *productRepository) GetBySlug(c context.Context, slug string) (*domain.Product, error) {
	q := r.q(c)

	p, err := q.GetProductBySlug(c, slug)
	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Product.GetBySlug slug=%s: %w", slug, err)
	}

	return r.assemble(c, q, p)
}

func (r *productRepository) RecalcMinPrice(ctx context.Context, productID int64) error {
	if err := r.q(ctx).RecalcProductMinPrice(ctx, productID); err != nil {
		return fmt.Errorf("catalog.repo.Product.RecalcMinPrice id=%d: %w", productID, err)
	}
	return nil
}

func (r *productRepository) assemble(c context.Context, q *db.Queries, row db.CatalogProduct) (*domain.Product, error) {
	id := row.ID
	p := productToDomain(row)

	variantRows, err := q.ListVariantsByProduct(c, id)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.GetByID variants id=%d: %w", id, err)
	}

	p.Variants = make([]domain.Variant, len(variantRows))
	idx := make(map[int64]int, len(variantRows))
	for i, val := range variantRows {
		p.Variants[i] = variantToDomain(val)
		idx[val.ID] = i
	}

	optRows, err := q.ListVariantOptionsByProduct(c, id)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.GetByID options id=%d: %w", id, err)
	}

	for _, o := range optRows {
		i := idx[o.ProductVariantID]
		p.Variants[i].Options = append(p.Variants[i].Options, domain.VariantOption{
			TypeID:  o.TypeID,
			Type:    o.TypeName,
			ValueID: o.ValueID,
			Value:   o.Value,
		})
	}

	imgRows, err := q.ListImagesByProduct(c, id)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.GetByID images id=%d: %w", id, err)
	}

	p.Images = make([]domain.Image, len(imgRows))
	for i, val := range imgRows {
		p.Images[i] = imageToDomain(val)
	}

	attrRows, err := q.ListProductAttributes(c, id)
	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Product.assemble attributes id=%d: %w", id, err)
	}

	p.Attributes = make([]domain.ProductAttribute, len(attrRows))
	for i, a := range attrRows {
		p.Attributes[i] = domain.ProductAttribute{
			AttributeID: a.AttributeID, Code: a.Code, Name: a.Name,
			ValueType: domain.AttributeValueType(a.ValueType), Unit: a.Unit,
			DisplayOrder: a.DisplayOrder, ValueText: a.ValueText, ValueNumber: a.ValueNumber,
		}
	}

	return p, nil
}
