ALTER TABLE catalog.products
  ADD COLUMN min_price bigint CHECK (min_price >= 0);

UPDATE catalog.products p
SET min_price = (
  SELECT MIN(pv.price) FROM catalog.product_variants pv
  WHERE pv.product_id = p.id AND pv.is_active
);

CREATE INDEX idx_products_min_price ON catalog.products (min_price, id)
  WHERE is_active AND deleted_at IS NULL AND min_price IS NOT NULL;