grpc-gateway:
 - POST /v1/auth/login
 - POST /v1/auth/create_user
 - POST /v1/auth/answer_challenge

 - GET /v1/catalog/products/{product_id}
 - POST /v1/catalog/products
 - PATCH /v1/catalog/products/{product_id}
 - DELETE /v1/catalog/products/{product_id}
 - GET /v1/catalog/products/list

 - GET /v1/cart
 - POST /v1/cart/items
 - PUT /v1/cart/items/{item_id}
 - DELETE /v1/cart/items/{item_id}
 - POST /v1/cart/checkout

 - GET /v1/inventory/{product_id}/availability
 - POST /v1/inventory/reserve
 - POST /v1/inventory/release
 - POST /v1/inventory/commit
 - GET /v1/inventory/skus
 - POST /v1/inventory/adjust_stock
 - POST /v1/inventory/add_sku
 - POST /v1/inventory/add_warehouse

 - POST /v1/payment/create-payment-session

normal:
 - /v1/auth/oauth/{provider}
 - /v1/auth/oauth/{provider}/callback
