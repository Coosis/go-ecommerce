-- name: GetStockLevelsByProductAndSkus :many
SELECT sku_id, warehouse_id, on_hand, reserved, updated_at 
FROM stock_levels sl JOIN skus s ON sl.sku_id = s.id
WHERE s.product_id = $1 
  AND s.id = ANY(sqlc.arg(sku_ids)::UUID[]);

-- name: GetStockLevelsByProductAndSkusAndWarehouse :many
SELECT sku_id, warehouse_id, on_hand, reserved, updated_at 
FROM stock_levels sl JOIN skus s ON sl.sku_id = s.id
WHERE s.product_id = $1 
  AND s.id = ANY(sqlc.arg(sku_ids)::UUID[]) 
  AND sl.warehouse_id = $2;


-- name: ReserveStock :one
WITH bumped AS (
  UPDATE stock_levels sl SET
  sl.reserved = sl.reserved + $4,
  sl.updated_at = NOW()
  WHERE sl.sku_id = $1 
    AND sl.warehouse_id = $2 
    AND (sl.on_hand - sl.reserved) >= $4
  RETURNING *
) INSERT INTO reservations (sku_id, warehouse_id, order_id, qty, expires_at)
SELECT $1, $2, $3, $4, NOW() + sqlc.arg(exp_duration)::INTERVAL
FROM bumped
RETURNING *;

-- name: ReleaseReservation :one
WITH released AS (
  UPDATE reservations
  SET status = 'released', updated_at = NOW()
  WHERE id = $1 AND status = 'active'
  RETURNING *
) UPDATE stock_levels sl SET
sl.reserved = sl.reserved - released.qty,
sl.updated_at = NOW()
FROM released
WHERE sl.sku_id = released.sku_id 
  AND sl.warehouse_id = released.warehouse_id
RETURNING sl.*;

-- name: CommitReservation :one
WITH committed AS (
  UPDATE reservations
  SET status = 'committed', updated_at = NOW()
  WHERE id = $1 
    AND status = 'active'
    AND expires_at > NOW()
  RETURNING *
) UPDATE stock_levels sl SET
sl.on_hand = sl.on_hand - committed.qty,
sl.reserved = sl.reserved - committed.qty,
sl.updated_at = NOW()
FROM committed
WHERE committed.sku_id = sl.sku_id 
  AND committed.warehouse_id = sl.warehouse_id
  AND (sl.on_hand - committed.qty) >= 0
RETURNING sl.*;

-- name: ExpireReservations :many
WITH expired AS (
  UPDATE reservations
  SET status = 'expired', updated_at = NOW()
  WHERE expires_at < NOW() 
    AND status = 'active'
  RETURNING *
), agg AS (
  SELECT sku_id, warehouse_id, SUM(qty) AS qty
  FROM expired
  GROUP BY sku_id, warehouse_id
) UPDATE stock_levels sl SET
sl.reserved = sl.reserved - agg.qty,
sl.updated_at = NOW()
FROM agg
WHERE agg.sku_id = sl.sku_id 
  AND agg.warehouse_id = sl.warehouse_id
RETURNING sl.*;

-- name: AdjustStockLevel :one
WITH adjusted AS (
  INSERT INTO stock_adjustments (sku_id, warehouse_id, delta, reason, created_by)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
), ins AS (
  INSERT INTO stock_levels (sku_id, warehouse_id, on_hand, reserved, updated_at)
  VALUES ($1, $2, GREATEST($3, 0), 0, NOW())
  ON CONFLICT (warehouse_id, sku_id)
  DO UPDATE SET 
    on_hand = GREATEST(stock_levels.on_hand + $3, 0), updated_at = NOW()
  WHERE (stock_levels.on_hand + $3) >= stock_levels.reserved
  RETURNING *
)
SELECT * FROM ins;

-- name: GetSkuCode :one
SELECT code FROM skus WHERE id = $1;

-- name: GetAllSkus :many
SELECT * FROM skus WHERE product_id = $1;

-- name: CreateSku :one
INSERT INTO skus (product_id, code) VALUES ($1, $2) RETURNING *;

-- name: CreateWarehouse :one
INSERT INTO warehouses (code, name) VALUES ($1, $2) RETURNING *;
