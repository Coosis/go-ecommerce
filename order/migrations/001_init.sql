-- +goose Up
CREATE TYPE order_status AS ENUM (
  'draft',           -- created but not paid yet
  'awaiting_payment',
  'paid',
  'canceled',
  'payment_failed',
  'fulfillment_started',
  'fulfilled',
  'refunded',
  'partially_refunded'
);

CREATE TABLE orders (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_number        TEXT UNIQUE NOT NULL,                     -- human-friendly; ULID/base36 OK

  user_id             UUID NOT NULL,
  email               TEXT,
  phone               TEXT,

  cart_id             UUID,                                     -- reference only (for forensics)
  currency            TEXT NOT NULL DEFAULT 'USD',
  subtotal_cents      BIGINT NOT NULL,
  discount_cents      BIGINT NOT NULL DEFAULT 0,
  shipping_cents      BIGINT NOT NULL DEFAULT 0,
  tax_cents           BIGINT NOT NULL DEFAULT 0,
  total_cents         BIGINT NOT NULL,
  status              order_status NOT NULL DEFAULT 'draft',

  payment_intent_id   TEXT,                                     -- from Payment Service/PSP

  shipping_address    TEXT NOT NULL,
  billing_address     TEXT,
  notes               TEXT,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  version             BIGINT NOT NULL DEFAULT 1
);

CREATE TABLE order_items (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id            UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id          UUID,          -- still store for BI; don't join for reads
  product_name        TEXT NOT NULL, -- snapshot for receipt
  sku_id              UUID,          -- immutable identifier used by inventory
  sku_code            TEXT,
  qty                 INT  NOT NULL CHECK (qty > 0),
  unit_price_cents    BIGINT NOT NULL,       -- snapshot at purchase time
  discount_cents      BIGINT NOT NULL DEFAULT 0,
  tax_rate_bp         INT NOT NULL DEFAULT 0, -- basis points (e.g., 825 = 8.25%)
  total_line_cents    BIGINT NOT NULL,       -- qty*(unit-price - discounts) + line-level tax if you do it that way
  price_version       BIGINT,                -- from Pricing/Catalog for later forensics
  metadata            TEXT                  -- arbitrary snapshot (attributes, options)
);

CREATE INDEX ON order_items(order_id);

-- +goose Down
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
