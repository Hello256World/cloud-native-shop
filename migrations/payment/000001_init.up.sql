CREATE SCHEMA IF NOT EXISTS payment;

CREATE TYPE payment.payment_status AS ENUM ('pending', 'succeeded', 'failed', 'refunded');

CREATE TABLE payment.payments (
  id                         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id                   bigint      NOT NULL,   -- reference to orders (no FK)
  user_id                    bigint      NOT NULL,   -- reference to identity (no FK)
  status                     payment.payment_status NOT NULL DEFAULT 'pending',
  amount                     bigint      NOT NULL CHECK (amount >= 0),
  gateway                    text        NOT NULL,
  type                       text        NOT NULL,   -- online | wallet | ...
  device                     text,
  retrieval_reference_number text        UNIQUE,
  idempotency_key            text        UNIQUE,
  failure_cause              text,
  description                text,
  created_at                 timestamptz NOT NULL DEFAULT now(),
  updated_at                 timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_payments_order ON payment.payments (order_id);
CREATE INDEX idx_payments_user  ON payment.payments (user_id, created_at DESC);

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON payment.payments
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();