# Products:
Get a product
GET `/products/{product_id}`

Create a product
POST `/products`
Body: `{ name: "shirt a", slug: "a-39-red", description: "A shirt", price_cents: 324 }`

Update an existing product
PATCH `/products/{product_id}`
Body: `{ name: "shirt a", slug: "a-39-red", description: "A shirt", price_cents: 324 }`

Delete a product
DELETE `/products/{productId}`

List products
GET `/products?category=&min_price=&max_price=&page=&per_page=`
