-- +goose Up
CREATE TABLE products (
  id           UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  name         TEXT    NOT NULL,
  slug         TEXT    NOT NULL UNIQUE,
  description  TEXT,
  price_cents  INTEGER NOT NULL,
  price_version BIGINT NOT NULL DEFAULT 1,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
  id         UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT    NOT NULL UNIQUE,
  slug       TEXT    NOT NULL UNIQUE,
  parent_id  UUID    REFERENCES categories(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE product_categories (
  product_id  UUID REFERENCES products(id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, category_id)
);

-- +goose Down
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;
