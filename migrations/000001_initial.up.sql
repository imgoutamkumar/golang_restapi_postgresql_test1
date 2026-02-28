-- Up migration: create initial tables for production
-- Enable pgcrypto extension for UUIDs (optional, if using UUIDs)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- Create ENUM type for order status
DO $$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'order_status'
) THEN CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'cancelled');
END IF;
END $$;
-- Create ENUM type for product status
DO $$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'product_status'
) THEN CREATE TYPE product_status AS ENUM (
    'draft',
    'active',
    'inactive',
    'archived'
);
END IF;
END $$;
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE
);
-- Insert the default role so it exists
INSERT INTO roles (name)
VALUES ('user') ON CONFLICT DO NOTHING;
-- Users table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL
);
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    password VARCHAR(255) NOT NULL,
    role_id UUID NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);
-- for multiple role --
-- CREATE TABLE IF NOT EXISTS user_roles (
--     user_id UUID NOT NULL,
--     role_id UUID NOT NULL,
--     PRIMARY KEY (user_id, role_id),
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
--     FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
-- );
CREATE TABLE IF NOT EXISTS user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    -- Contact info specific to this address (Crucial for gifting/office deliveries)
    receiver_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    -- Address breakdown
    address_line_1 VARCHAR(255) NOT NULL,
    -- Street, House No.
    address_line_2 VARCHAR(255),
    -- Apartment, Suite, Landmark
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'India',
    postal_code VARCHAR(20) NOT NULL,
    -- Metadata
    address_type VARCHAR(20) NOT NULL CHECK (address_type IN ('home', 'work', 'other')),
    is_default BOOLEAN NOT NULL DEFAULT false,
    -- Geospatial (Optional but recommended for delivery routing/ETA)
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    -- Audit fields matching your users table
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_addresses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
      type VARCHAR(50) NOT NULL, -- hero, sale, promo,
    is_active BOOLEAN DEFAULT true,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ ,
    CONSTRAINT fk_banners_users FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS banner_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    banner_id UUID NOT NULL,

    image_url TEXT NOT NULL,
    public_id VARCHAR(255) NOT NULL,

    link_url TEXT,
    title VARCHAR(255),
    description TEXT,

    is_active BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_banner
        FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,

    CONSTRAINT fk_banner_user
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_category UNIQUE(name, parent_id)
);
CREATE TABLE attribute_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);
CREATE TABLE attribute_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attribute_type_id UUID REFERENCES attribute_types(id) ON DELETE CASCADE,
    value VARCHAR(100) NOT NULL,
    meta_info VARCHAR(100),
    UNIQUE(attribute_type_id, value)
);
-- Products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    brand_id UUID NOT NULL,
    category_id UUID NOT NULL,
    short_description VARCHAR(500) NOT NULL,
    description TEXT,
    base_price NUMERIC(10, 2) NOT NULL,
    currency CHAR(3) DEFAULT 'INR',
    status product_status NOT NULL DEFAULT 'draft',
    is_returnable BOOLEAN DEFAULT true,
    is_cod_available BOOLEAN DEFAULT true,
    average_rating NUMERIC(3, 2) DEFAULT 0,
    rating_count INT DEFAULT 0,
    created_by UUID NOT NULL,
    details JSONB DEFAULT '{}'::jsonb,
    slug VARCHAR(160) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_products_users FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_products_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE RESTRICT,
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) UNIQUE NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    discount_percent NUMERIC(10, 2) CHECK (
        discount_percent >= 0
        AND discount_percent <= 100
    ),
    is_default BOOLEAN DEFAULT false,
    is_wishlisted BOOLEAN DEFAULT false,
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    status VARCHAR(20) DEFAULT 'active',
    slug VARCHAR(120) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE variant_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    attribute_value_id UUID REFERENCES attribute_values(id) ON DELETE CASCADE,
    UNIQUE(variant_id, attribute_value_id)
);
CREATE TABLE IF NOT EXISTS product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID NOT NULL,
    image_url TEXT NOT NULL,
    public_id VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_product_images_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    parent_id UUID NULL,
    rating SMALLINT CHECK (
        rating BETWEEN 1 AND 5
    ),
    comment_text TEXT,
    is_edited BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);
-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    order_number VARCHAR(30) UNIQUE NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    subtotal NUMERIC(10, 2) NOT NULL,
    discount_amount NUMERIC(10, 2) DEFAULT 0,
    tax_amount NUMERIC(10, 2) DEFAULT 0,
    shipping_amount NUMERIC(10, 2) DEFAULT 0,
    total_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    product_name VARCHAR(150) NOT NULL,
    product_variant_name VARCHAR(150) NOT NULL,
    product_variant_price NUMERIC(10, 2) NOT NULL,
    discount_percent NUMERIC(10, 2) DEFAULT 0,
    quantity INT NOT NULL CHECK (quantity > 0),
    total_price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id),
    CONSTRAINT uq_order_variant UNIQUE (order_id, variant_id)
);
-- Cart table
-- carts (1 cart per user)
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- cart items (many products per cart)
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id),
    CONSTRAINT uq_cart_product UNIQUE (cart_id, variant_id)
);
-- 1. The Wishlist Container
CREATE TABLE wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    -- In production, don't use a strict FK if users are in a different microservice
    name VARCHAR(100) DEFAULT 'My Wishlist',
    is_public BOOLEAN DEFAULT false,
    -- Allows sharing
    share_token VARCHAR(64) UNIQUE,
    -- e.g., "abc123xyz" for myntra.com/wishlist/abc123xyz
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE wishlist_items (
    id BIGSERIAL PRIMARY KEY,
    wishlist_id UUID REFERENCES wishlists(id) ON DELETE CASCADE,
    -- We store BOTH product_id and variant_id. 
    -- variant_id tells us exactly what they want (Size XL, Red).
    -- product_id is denormalized here to make "Is this product in my wishlist?" queries lightning fast.
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    -- Crucial FAANG Feature: Track the price when they added it
    added_price DECIMAL(12, 2) NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (wishlist_id, variant_id)
);
CREATE TABLE IF NOT EXISTS password_reset (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    otp_hash VARCHAR(255) NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ NOT NULL,
    locked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_password_reset_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- 1. Warehouses
-- Stores physical locations where inventory is kept.
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    location_address TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_warehouses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- 2. Inventory Levels (The Current Snapshot)
-- Tracks the exact current count of a specific variant in a specific warehouse.
CREATE TABLE IF NOT EXISTS inventory_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
    -- Stock ready to be sold
    available_stock INT NOT NULL DEFAULT 0 CHECK (available_stock >= 0),
    -- Stock currently in someone's checkout process, but not yet paid for
    reserved_stock INT NOT NULL DEFAULT 0 CHECK (reserved_stock >= 0),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Ensure we only have one row per variant per warehouse
    UNIQUE (variant_id, warehouse_id)
);
-- 3. Inventory Ledgers (The Immutable Audit Trail)
-- An append-only log of EVERY time stock goes up or down. 
-- NEVER UPDATE rows in this table, only INSERT.
CREATE TABLE IF NOT EXISTS inventory_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
    -- Positive numbers (+50) for restock/returns. Negative numbers (-1) for sales/damage.
    quantity_change INT NOT NULL,
    -- Why did the stock change? 
    -- e.g., 'RESTOCK', 'ORDER_PLACED', 'RETURNED_TO_STOCK', 'DAMAGED_IN_TRANSIT'
    reason VARCHAR(50) NOT NULL,
    -- Links to an Order ID or a Purchase Order ID for perfect traceability
    reference_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- 1. Standard lookup index (Used every time a user opens their "My Addresses" page)
CREATE INDEX IF NOT EXISTS idx_user_addresses_user_id ON user_addresses(user_id)
WHERE deleted_at IS NULL;
-- 2. The FAANG "Single Default" Constraint (Partial Unique Index)
-- This guarantees that a user can ONLY have ONE default address at a time.
-- If you try to set two addresses to is_default=true, PostgreSQL will block it.
CREATE UNIQUE INDEX IF NOT EXISTS idx_single_default_address ON user_addresses(user_id)
WHERE is_default = true
    AND deleted_at IS NULL;
CREATE INDEX idx_comments_product_id ON comments(product_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_not_deleted ON comments(product_id)
WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_banner_active ON banners(is_active);
CREATE INDEX idx_banner_images_banner_id ON banner_images(banner_id);
-- Index for fast lookups when querying a product's stock
CREATE INDEX idx_inventory_levels_variant ON inventory_levels(variant_id);
-- Index for auditing and calculating historical stock
CREATE INDEX idx_inventory_ledgers_variant ON inventory_ledgers(variant_id);
CREATE INDEX idx_inventory_ledgers_created_at ON inventory_ledgers(created_at);
-- Ensure only one primary image per product
CREATE UNIQUE INDEX IF NOT EXISTS ux_product_primary_image ON product_images (variant_id)
WHERE is_primary = true;
-- Trigger functions to auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Attach triggers
CREATE TRIGGER users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER products_updated_at BEFORE
UPDATE ON products FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER orders_updated_at BEFORE
UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER carts_updated_at BEFORE
UPDATE ON carts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER product_images_updated_at BEFORE
UPDATE ON product_images FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- Additional Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_products_brand ON products(brand_id);
CREATE INDEX idx_products_created_by ON products(created_by);
CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_variant_attributes_variant ON variant_attributes(variant_id);
CREATE INDEX idx_variant_attributes_value ON variant_attributes(attribute_value_id);