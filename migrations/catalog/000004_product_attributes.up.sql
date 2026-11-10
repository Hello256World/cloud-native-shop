CREATE TABLE catalog.attribute_definitions (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  code          text        NOT NULL,
  name          text        NOT NULL,
  value_type    text        NOT NULL,
  unit          text,
  is_filterable boolean     NOT NULL DEFAULT false,
  display_order integer     NOT NULL DEFAULT 0,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_attribute_definitions_code UNIQUE (code),
  CONSTRAINT ck_attribute_definitions_type CHECK (value_type IN ('text', 'number'))
);

CREATE TABLE catalog.product_attributes (
  product_id   bigint  NOT NULL REFERENCES catalog.products (id) ON DELETE CASCADE,
  attribute_id bigint  NOT NULL REFERENCES catalog.attribute_definitions (id),
  value_text   text,
  value_number numeric,
  PRIMARY KEY (product_id, attribute_id),
  CONSTRAINT ck_product_attributes_one_value
    CHECK ((value_text IS NOT NULL) <> (value_number IS NOT NULL))
);
CREATE INDEX idx_pa_number ON catalog.product_attributes (attribute_id, value_number)
  WHERE value_number IS NOT NULL;
CREATE INDEX idx_pa_text ON catalog.product_attributes (attribute_id, value_text)
  WHERE value_text IS NOT NULL;

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON catalog.attribute_definitions
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();