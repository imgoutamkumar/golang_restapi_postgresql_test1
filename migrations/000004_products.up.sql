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
CREATE TABLE IF NOT EXISTS banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    -- hero, sale, promo,
    is_active BOOLEAN DEFAULT true,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
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
    CONSTRAINT fk_banner FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    CONSTRAINT fk_banner_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_category UNIQUE(name, parent_id)
);
CREATE TABLE IF NOT EXISTS attribute_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS attribute_values (
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
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) UNIQUE NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    discount_percent NUMERIC(10, 2) CHECK (
        discount_percent >= 0
        AND discount_percent <= 100
    ),
    is_default BOOLEAN DEFAULT false,
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    status VARCHAR(20) DEFAULT 'active',
    slug VARCHAR(120) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE IF NOT EXISTS variant_attributes (
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

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    rating INT NOT NULL CHECK (
        rating BETWEEN 1 AND 5
    ),
    review TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(user_id, product_id)
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    review_id UUID NULL,
    -- optional link to a review
    parent_id UUID NULL,
    -- for nesting replies
    comment_text TEXT NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_review FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_products_slug_active ON products(slug)
WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX ux_default_variant_per_product ON product_variants(product_id)
WHERE is_default = true;
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_product_id ON comments(product_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_not_deleted ON comments(product_id)
WHERE deleted_at IS NULL;
CREATE INDEX idx_banner_active ON banners(is_active);
CREATE INDEX idx_banner_images_banner_id ON banner_images(banner_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_product_primary_image ON product_images (variant_id)
WHERE is_primary = true;
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_products_brand ON products(brand_id);
CREATE INDEX idx_products_created_by ON products(created_by);
CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_variant_attributes_variant ON variant_attributes(variant_id);
CREATE INDEX idx_variant_attributes_value ON variant_attributes(attribute_value_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_variants_price ON product_variants(price);
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_comments_review_id ON comments(review_id);