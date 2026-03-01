-- Cart table
-- carts (1 cart per user)
CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- cart items (many products per cart)
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    CONSTRAINT uq_cart_product UNIQUE (cart_id, variant_id)
);
-- 1. The Wishlist Container
CREATE TABLE IF NOT EXISTS wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    -- In production, don't use a strict FK if users are in a different microservice
    name VARCHAR(100) DEFAULT 'My Wishlist',
    is_public BOOLEAN DEFAULT false,
    -- Allows sharing
    share_token VARCHAR(64) UNIQUE,
    -- e.g., "abc123xyz" for myntra.com/wishlist/abc123xyz
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wishlist_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlist_id UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    -- We store BOTH product_id and variant_id. 
    -- variant_id tells us exactly what they want (Size XL, Red).
    -- product_id is denormalized here to make "Is this product in my wishlist?" queries lightning fast.
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    -- Crucial FAANG Feature: Track the price when they added it
    added_price DECIMAL(12, 2) NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (wishlist_id, variant_id)
);
CREATE INDEX IF NOT EXISTS idx_wishlist_user ON wishlists(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_variant_id ON cart_items(variant_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_items_product_id ON wishlist_items(product_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_items_variant_id ON wishlist_items(variant_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_wishlist_share_token ON wishlists(share_token)
WHERE share_token IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS ux_carts_user_active ON carts(user_id)
WHERE deleted_at IS NULL;