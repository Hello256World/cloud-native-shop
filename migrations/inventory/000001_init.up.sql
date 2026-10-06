CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TYPE inventory.reservation_status AS ENUM ('active', 'committed', 'released', 'expired');

CREATE TABLE inventory.stock_items (
  product_variant_id bigint      PRIMARY KEY,
  quantity_on_hand   integer     NOT NULL DEFAULT 0 CHECK (quantity_on_hand >= 0),
  quantity_reserved  integer     NOT NULL DEFAULT 0 CHECK (quantity_reserved >= 0),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  CHECK (quantity_reserved <= quantity_on_hand)
);

CREATE TABLE inventory.stock_reservations (
  id                 bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id           bigint      NOT NULL,
  product_variant_id bigint      NOT NULL REFERENCES inventory.stock_items (product_variant_id),
  quantity           integer     NOT NULL CHECK (quantity > 0),
  status             inventory.reservation_status NOT NULL DEFAULT 'active',
  expires_at         timestamptz NOT NULL,
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_reservations_order  ON inventory.stock_reservations (order_id);
CREATE INDEX idx_reservations_expiry ON inventory.stock_reservations (expires_at) WHERE status = 'active';

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON inventory.stock_items
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON inventory.stock_reservations
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();