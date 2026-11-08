CREATE TABLE catalog.compare_products (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  product_id bigint      NOT NULL REFERENCES catalog.products (id) ON DELETE CASCADE,
  name       text        NOT NULL,
  link       text        NOT NULL,
  price      bigint      NOT NULL CHECK (price >= 0),
  image      text,
  is_active  boolean     NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_compare_products_product ON catalog.compare_products (product_id);
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.compare_products
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();