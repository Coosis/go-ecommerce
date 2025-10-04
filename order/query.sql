-- name: CreateOrder :one
WITH new_order AS (
  INSERT INTO orders (
    order_number, 

    user_id, email, phone, 

    cart_id, currency,
    subtotal_cents, discount_cents, shipping_cents, tax_cents, total_cents,

    payment_intent_id, 

    shipping_address, billing_address, notes
  ) VALUES (
    $1, 
    $2, $3, $4,
    $5, $6,
    $7, $8, $9, $10, $11,
    $12, 
    sqlc.arg(shipping_addr), sqlc.arg(billing_addr), $13
  ) RETURNING *
) SELECT * FROM new_order;

-- name: CancelOrder :one
WITH canceled AS (
  UPDATE orders SET
    status = 'canceled',
    updated_at = NOW(),
    version = version + 1
  WHERE id = $1 
    AND status IN ('draft', 'awaiting_payment', 'payment_failed')
  RETURNING *
) SELECT * FROM canceled;

-- name: AddOrderItem :one
INSERT INTO order_items (
  order_id, 
  product_id, product_name,
  sku_id, sku_code, 
  qty, 
  unit_price_cents, discount_cents, tax_rate_bp, total_line_cents,
  price_version, metadata
) VALUES (
  $1, 
  $2, $3,
  $4, $5, 
  $6, 
  $7, $8, $9, $10,
  $11, sqlc.arg(metadata)
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;
