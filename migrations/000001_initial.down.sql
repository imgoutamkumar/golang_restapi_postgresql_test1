-- =========================
-- DROP TRIGGERS
-- =========================
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP TRIGGER IF EXISTS products_updated_at ON products;
DROP TRIGGER IF EXISTS orders_updated_at ON orders;
DROP TRIGGER IF EXISTS carts_updated_at ON carts;
DROP TRIGGER IF EXISTS product_images_updated_at ON product_images;

-- =========================
-- DROP TRIGGER FUNCTION
-- =========================
DROP FUNCTION IF EXISTS set_updated_at();

-- =========================
-- DROP INDEXES
-- =========================
DROP INDEX IF EXISTS ux_product_primary_image;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_users_email;

-- =========================
-- DROP ECOM TABLES (child → parent)
-- =========================

-- cart + order
DROP TABLE IF EXISTS cart_items CASCADE;
DROP TABLE IF EXISTS carts CASCADE;

DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;

-- password reset
DROP TABLE IF EXISTS password_reset CASCADE;

-- variant system
DROP TABLE IF EXISTS variant_attributes CASCADE;
DROP TABLE IF EXISTS product_variants CASCADE;

-- product images
DROP TABLE IF EXISTS product_images CASCADE;

-- products
DROP TABLE IF EXISTS products CASCADE;

-- catalog
DROP TABLE IF EXISTS attribute_values CASCADE;
DROP TABLE IF EXISTS attribute_types CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS brands CASCADE;

-- RBAC
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;

-- users
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- =========================
-- DROP ENUM TYPES
-- =========================
DROP TYPE IF EXISTS product_status CASCADE;
DROP TYPE IF EXISTS order_status CASCADE;
