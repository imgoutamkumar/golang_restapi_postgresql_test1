-- =========================
-- DOWN MIGRATION: TRIGGERS
-- =========================

-- Drop triggers first
DROP TRIGGER IF EXISTS orders_updated_at ON orders;
DROP TRIGGER IF EXISTS wishlists_updated_at ON wishlists;
DROP TRIGGER IF EXISTS carts_updated_at ON carts;
DROP TRIGGER IF EXISTS product_images_updated_at ON product_images;
DROP TRIGGER IF EXISTS products_updated_at ON products;
DROP TRIGGER IF EXISTS users_updated_at ON users;

-- Drop the trigger function
DROP FUNCTION IF EXISTS set_updated_at();