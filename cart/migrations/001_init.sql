-- +goose Up
CREATE TABLE carts (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID,
  version      BIGINT NOT NULL DEFAULT 1,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  status       TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','merged','abandoned','checked_out'))
);

create table cart_items (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id              UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  product_id           UUID NOT NULL,
  sku_id               UUID NOT NULL,
  qty                  INTEGER NOT NULL CHECK (qty > 0),
  price_cents_snapshot INTEGER NOT NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  unique (cart_id, product_id)
);

CREATE UNIQUE INDEX uniq_active_cart_per_user
  ON carts(user_id)
  WHERE status = 'active';

-- +goose Down
DROP INDEX IF EXISTS uniq_active_cart_per_user;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
