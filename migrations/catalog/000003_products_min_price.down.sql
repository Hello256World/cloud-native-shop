DROP INDEX IF EXISTS catalog.idx_products_min_price;
ALTER TABLE catalog.products DROP COLUMN IF EXISTS min_price;