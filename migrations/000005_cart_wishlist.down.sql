DROP INDEX IF EXISTS ux_carts_user_active;
DROP INDEX IF EXISTS ux_wishlist_share_token;
DROP INDEX IF EXISTS idx_wishlist_items_variant_id;
DROP INDEX IF EXISTS idx_wishlist_items_product_id;
DROP INDEX IF EXISTS idx_wishlist_items_wishlist_id;
DROP INDEX IF EXISTS idx_cart_items_variant_id;
DROP INDEX IF EXISTS idx_cart_items_cart_id;
DROP INDEX IF EXISTS idx_wishlist_user;

DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS carts;