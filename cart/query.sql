-- name: CreateCartById :one
INSERT INTO carts(id) VALUES ($1) RETURNING *;

-- name: CreateCartByIdAndUserId :one
INSERT INTO carts(id, user_id) VALUES ($1, $2) RETURNING *;

-- name: AssociateCartWithUser :exec
UPDATE carts SET user_id = $2, updated_at = NOW() WHERE id = $1;

-- name: FindCartByUserId :one
SELECT * FROM carts WHERE user_id = $1;

-- name: InsertItemToCart :one
INSERT INTO cart_items (
  cart_id,
  product_id,
  sku_id,
  qty,
  price_cents_snapshot
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: RemoveItemFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;

-- name: GetCart :one
SELECT * FROM carts WHERE id = $1;

-- name: GetCartItems :many
SELECT * FROM cart_items WHERE cart_id = $1;

-- name: GetActiveCartByUserId :one
INSERT INTO carts (user_id, status)
VALUES ($1, 'active')
ON CONFLICT (user_id) WHERE status = 'active'
DO UPDATE SET updated_at = NOW()
RETURNING *;

-- name: GetActiveCartByCartId :one
INSERT INTO carts (id, status)
VALUES ($1, 'active')
ON CONFLICT (id)
DO UPDATE SET updated_at = NOW()
RETURNING *;

-- *** name: RotateAndMergeCartForUser :one
-- WITH old_active AS (
--   SELECT id
--   FROM carts c
--   WHERE c.user_id = $1 AND c.status = 'active'
--   FOR UPDATE
-- ),
-- also_active AS (
--   SELECT id
--   FROM carts c
--   WHERE c.id = $2 AND c.status = 'active'
--   FOR UPDATE
-- ),
-- merged AS (
--   UPDATE carts
--   SET status = 'merged', updated_at = NOW()
--   WHERE id IN (
--     SELECT id FROM old_active
--     UNION
--     SELECT id FROM also_active
--   )
--   RETURNING id
-- ),
-- new_cart AS (
--   INSERT INTO carts (id, user_id, status)
--   VALUES (gen_random_uuid(), $1, 'active')
--   RETURNING id, user_id, version, created_at, updated_at, status
-- ),
-- moved AS (
--   INSERT INTO cart_items (cart_id, product_id, qty, price_cents_snapshot)
--   SELECT (SELECT id FROM new_cart), ci.product_id, ci.qty, ci.price_cents_snapshot
--   FROM cart_items ci
--   WHERE ci.cart_id IN (SELECT id FROM merged)
--   ON CONFLICT (cart_id, product_id)
--   DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty,
--                 updated_at = NOW()
--   RETURNING 1
-- )
-- SELECT * FROM new_cart;

-- name: RotateAndMergeCartForUser :one
WITH upsert AS (
  INSERT INTO carts (id, user_id, status)
  VALUES (gen_random_uuid(), @user_id, 'active')
  ON CONFLICT (user_id) WHERE status = 'active'
  DO UPDATE SET updated_at = NOW()
  RETURNING *
),
to_merge AS (
  UPDATE carts c
  SET status = 'merged', updated_at = NOW()
  WHERE @user_id IS NOT NULL
    AND c.id = @old_cart_id
    AND c.status = 'active'
    AND c.id <> (SELECT id FROM upsert)
  RETURNING c.id
),
moved AS (
  INSERT INTO cart_items (cart_id, product_id, qty, price_cents_snapshot)
  SELECT (SELECT id FROM upsert), ci.product_id, ci.qty, ci.price_cents_snapshot
  FROM cart_items ci
  WHERE ci.cart_id IN (SELECT id FROM to_merge)
  ON CONFLICT (cart_id, product_id)
  DO UPDATE SET
    qty = cart_items.qty + EXCLUDED.qty,
    updated_at = NOW()
  RETURNING 1
)
SELECT
  c.id, c.user_id, c.version, c.created_at, c.updated_at, c.status
FROM upsert c;
