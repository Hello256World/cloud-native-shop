CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE catalog.categories (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  parent_id  bigint      REFERENCES catalog.categories (id),
  name       text        NOT NULL,
  slug       text        NOT NULL UNIQUE,
  image_name text,
  is_active  boolean     NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);
CREATE INDEX idx_categories_parent ON catalog.categories (parent_id);

CREATE TABLE catalog.products (
  id                bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  category_id       bigint      NOT NULL REFERENCES catalog.categories (id),
  name              text        NOT NULL,
  slug              text        NOT NULL UNIQUE,
  description       text,
  thumbnail         text        NOT NULL,
  shipment_weight_g integer     NOT NULL CHECK (shipment_weight_g >= 0),
  is_active         boolean     NOT NULL DEFAULT true,
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  deleted_at        timestamptz
);
CREATE INDEX idx_products_category ON catalog.products (category_id) WHERE deleted_at IS NULL;

CREATE TABLE catalog.product_images (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  product_id bigint      NOT NULL REFERENCES catalog.products (id) ON DELETE CASCADE,
  image_name text        NOT NULL,
  priority   integer     NOT NULL DEFAULT 0,
  is_active  boolean     NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_product_images_product ON catalog.product_images (product_id, priority);

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

CREATE TABLE catalog.variant_types (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name          text        NOT NULL UNIQUE,
  data_type     text        NOT NULL,
  display_order integer     NOT NULL DEFAULT 0,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE catalog.variant_values (
  id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  variant_type_id bigint      NOT NULL REFERENCES catalog.variant_types (id) ON DELETE CASCADE,
  value           text        NOT NULL,
  display_order   integer     NOT NULL DEFAULT 0,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (variant_type_id, value)
);

CREATE TABLE catalog.product_variants (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  product_id bigint      NOT NULL REFERENCES catalog.products (id) ON DELETE CASCADE,
  sku        text        NOT NULL UNIQUE,
  price      bigint      NOT NULL CHECK (price >= 0),
  is_active  boolean     NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_product_variants_product ON catalog.product_variants (product_id);

CREATE TABLE catalog.product_variant_values (
  product_variant_id bigint NOT NULL REFERENCES catalog.product_variants (id) ON DELETE CASCADE,
  variant_value_id   bigint NOT NULL REFERENCES catalog.variant_values (id),
  PRIMARY KEY (product_variant_id, variant_value_id)
);
CREATE INDEX idx_pvv_value ON catalog.product_variant_values (variant_value_id);

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.categories
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.products
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.product_images
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.compare_products
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.variant_types
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.variant_values
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.product_variants
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();