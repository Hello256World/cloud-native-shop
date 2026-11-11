package postgres

import (
	"context"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

func (r *productRepository) ReplaceAttributes(ctx context.Context, productID int64, attrs []domain.ProductAttribute) error {
	q := r.q(ctx)

	if err := q.DeleteProductAttributes(ctx, productID); err != nil {
		return fmt.Errorf("catalog.repo.Product.ReplaceAttributes delete product=%d: %w", productID, err)
	}

	for _, a := range attrs {
		if err := q.InsertProductAttribute(ctx, db.InsertProductAttributeParams{
			ProductID:   productID,
			AttributeID: a.AttributeID,
			ValueText:   a.ValueText,
			ValueNumber: a.ValueNumber,
		}); err != nil {
			return fmt.Errorf("catalog.repo.Product.ReplaceAttributes insert product=%d attr=%d: %w", productID, a.AttributeID, err)
		}
	}
	return nil
}
