-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, slug, description, price_cents)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET 
	name = $2,
	slug = $3,
	description = $4,
	price_cents = $5,
	updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT p.* FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
WHERE
	(sqlc.narg('max_price_cents')::int IS NULL OR p.price_cents <= sqlc.narg('max_price_cents'))
	AND (sqlc.narg('min_price_cents')::int IS NULL OR price_cents >= sqlc.narg('min_price_cents'))
	AND (c.slug = sqlc.narg('slug') OR sqlc.narg('slug')::text IS NULL)
ORDER BY p.name
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);
