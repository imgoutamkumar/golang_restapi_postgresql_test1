-- Down migration: rollback catalog, banners, and related entities

-- 1. Drop indexes (optional, PostgreSQL auto-drops with tables, but safe)
DROP INDEX IF EXISTS idx_comments_review_id;
DROP INDEX IF EXISTS idx_reviews_product_id;
DROP INDEX IF EXISTS idx_variants_price;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_variant_attributes_value;
DROP INDEX IF EXISTS idx_variant_attributes_variant;
DROP INDEX IF EXISTS idx_variants_product;
DROP INDEX IF EXISTS idx_products_created_by;
DROP INDEX IF EXISTS idx_products_brand;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_categories_parent;
DROP INDEX IF EXISTS ux_product_primary_image;
DROP INDEX IF EXISTS idx_banner_images_banner_id;
DROP INDEX IF EXISTS idx_banner_active;
DROP INDEX IF EXISTS idx_comments_not_deleted;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_product_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS ux_default_variant_per_product;
DROP INDEX IF EXISTS ux_products_slug_active;

-- 2. Drop deepest child tables first
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS reviews;

DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS variant_attributes;

DROP TABLE IF EXISTS product_variants;

-- 3. Drop main product table
DROP TABLE IF EXISTS products;

-- 4. Drop attribute system
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attribute_types;

-- 5. Drop category & brand
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS brands;

-- 6. Drop banners
DROP TABLE IF EXISTS banner_images;
DROP TABLE IF EXISTS banners;

-- 7. Drop ENUM (last, after all dependencies are gone)
DROP TYPE IF EXISTS product_status;