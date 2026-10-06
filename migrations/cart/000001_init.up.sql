CREATE SCHEMA IF NOT EXISTS cart;

CREATE TABLE cart.carts (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id    bigint      NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE cart.cart_items (
  id                 bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  cart_id            bigint      NOT NULL REFERENCES cart.carts (id) ON DELETE CASCADE,
  product_variant_id bigint      NOT NULL,
  quantity           integer     NOT NULL CHECK (quantity > 0),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cart_id, product_variant_id)
);

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON cart.carts
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON cart.cart_items
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();