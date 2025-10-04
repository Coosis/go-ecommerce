Create & own a cart for a guest or authenticated user.
CRUD cart items (add, update qty, remove).
Persist price snapshots per line (unit_price_snapshot, currency) for price stability.
Compute lightweight totals: line totals + subtotal only.
Merge carts (guest → user) on login.
Validate items against Catalog (existence + purchasable flag).
Hold client intent: optional coupon code placeholder; optional notes.
Expose a clean read model for Checkout and UI.
Emit events (e.g., CartCreated, CartItemChanged, CartMerged, CartEmptied).
