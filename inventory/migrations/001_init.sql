-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE warehouses (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code        TEXT UNIQUE NOT NULL,
  name        TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE skus (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id  UUID NOT NULL,
  code        TEXT UNIQUE NOT NULL, -- e.g. "TSHIRT-BLK-M"
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stock_levels (
  warehouse_id UUID NOT NULL REFERENCES warehouses(id),
  sku_id       UUID NOT NULL REFERENCES skus(id),
  on_hand      INTEGER NOT NULL DEFAULT 0,     -- physical units
  reserved     INTEGER NOT NULL DEFAULT 0,     -- soft holds
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (warehouse_id, sku_id),
  CHECK (on_hand >= 0),
  CHECK (reserved >= 0)
);

-- A reservation is a temporary hold while checkout is in progress.
CREATE TABLE reservations (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sku_id        UUID NOT NULL REFERENCES skus(id),
  warehouse_id  UUID NOT NULL REFERENCES warehouses(id),
  order_id        UUID NOT NULL, -- your order_id
  qty           INTEGER NOT NULL CHECK (qty > 0),
  status        TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','released','committed','expired')),
  expires_at    TIMESTAMPTZ NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (sku_id, warehouse_id, order_id)
);

-- Immutable journal for adjustments & audit trail.
CREATE TABLE stock_adjustments (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sku_id        UUID NOT NULL REFERENCES skus(id),
  warehouse_id  UUID NOT NULL REFERENCES warehouses(id),
  delta         INTEGER NOT NULL, -- +receiving, -damage, etc.
  reason        TEXT NOT NULL,    -- 'receiving','return','correction',...
  created_by    TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Outbox for reliable events (use in the same tx).
CREATE TABLE outbox (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  topic        TEXT NOT NULL,
  payload      JSONB NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  delivered_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS stock_adjustments;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS stock_levels;
DROP TABLE IF EXISTS skus;
DROP TABLE IF EXISTS warehouses;
