CREATE SCHEMA IF NOT EXISTS orders;

CREATE TYPE orders.order_status AS ENUM (
  'pending',
  'awaiting_payment',
  'paid',
  'processing',
  'shipped',
  'delivered',
  'cancelled',
  'rejected'
);

CREATE TABLE orders.orders (
  id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id          bigint      NOT NULL,
  status           orders.order_status NOT NULL DEFAULT 'pending',

  customer_name    text        NOT NULL,
  phone            text        NOT NULL,
  delivery_address text        NOT NULL,
  delivery_method  text        NOT NULL,

  description      text,
  rejection_reason text,
  total_weight_g   integer     NOT NULL DEFAULT 0 CHECK (total_weight_g >= 0),
  subtotal_amount  bigint      NOT NULL CHECK (subtotal_amount >= 0),
  shipping_amount  bigint      NOT NULL DEFAULT 0 CHECK (shipping_amount >= 0),
  total_amount     bigint      NOT NULL CHECK (total_amount >= 0),

  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_orders_user   ON orders.orders (user_id, created_at DESC);
CREATE INDEX idx_orders_status ON orders.orders (status);

CREATE TABLE orders.order_items (
  id                 bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id           bigint  NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
  product_variant_id bigint  NOT NULL,
  product_name       text    NOT NULL,   -- snapshot
  sku                text    NOT NULL,   -- snapshot
  unit_price         bigint  NOT NULL CHECK (unit_price >= 0),
  quantity           integer NOT NULL CHECK (quantity > 0)
);
CREATE INDEX idx_order_items_order ON orders.order_items (order_id);

CREATE TABLE orders.order_status_history (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id    bigint      NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
  from_status orders.order_status,
  to_status   orders.order_status NOT NULL,
  changed_by  bigint,
  note        text,
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_order_status_history_order ON orders.order_status_history (order_id, created_at);

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON orders.orders
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();