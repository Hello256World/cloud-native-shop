CREATE SCHEMA IF NOT EXISTS identity;

CREATE TYPE identity.gender AS ENUM ('male', 'female', 'other');

CREATE TABLE identity.roles (
  id   smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name text NOT NULL UNIQUE
);

INSERT INTO identity.roles (name) VALUES ('customer'), ('admin'), ('super_admin');

CREATE TABLE identity.users (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  full_name     text        NOT NULL,
  email         citext      UNIQUE,
  phone         text        NOT NULL UNIQUE,
  password_hash text,
  birthday      date,
  gender        identity.gender,
  is_active     boolean     NOT NULL DEFAULT true,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz
);

CREATE TABLE identity.user_roles (
  user_id bigint   NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  role_id smallint NOT NULL REFERENCES identity.roles (id),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE identity.addresses (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id       bigint      NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  receiver_name text        NOT NULL,
  phone         text        NOT NULL,
  address       text        NOT NULL,
  postal_code   text,
  no            text,
  unit          text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz
);
CREATE INDEX idx_addresses_user ON identity.addresses (user_id) WHERE deleted_at IS NULL;

CREATE TABLE identity.refresh_tokens (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id    bigint      NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  token_hash text        NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_refresh_tokens_user ON identity.refresh_tokens (user_id);

CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON identity.users
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();
CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON identity.addresses
  FOR EACH ROW EXECUTE FUNCTION platform.set_updated_at();