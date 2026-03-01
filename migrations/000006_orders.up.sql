-- Create ENUM type for order status
DO $$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'order_status'
) THEN CREATE TYPE order_status AS ENUM (
    'pending',
    'confirmed',
    'paid',
    'shipped',
    'delivered',
    'cancelled',
    'returned'
);
END IF;
END $$;
DO $$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'payment_status'
) THEN CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');
END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'shipment_status'
    ) THEN
        CREATE TYPE shipment_status AS ENUM (
            'pending','packed','shipped','out_for_delivery','delivered','failed'
        );
    END IF;
END $$;
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
    payment_status payment_status DEFAULT 'pending',
    payment_method VARCHAR(50),
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
    CONSTRAINT fk_order_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
        CONSTRAINT uq_order_variant UNIQUE (order_id, variant_id)
);
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',
    provider VARCHAR(50),
    -- razorpay, stripe
    provider_payment_id VARCHAR(255) NOT NULL UNIQUE,
    status payment_status DEFAULT 'pending',
    -- pending, success, failed
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    address_id UUID NOT NULL,
    courier_name VARCHAR(100),
    receiver_name VARCHAR(100),
    phone_number VARCHAR(20),
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    status shipment_status NOT NULL DEFAULT 'pending',
    -- pending, shipped, out_for_delivery, delivered
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT fk_shipment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_shipment_address FOREIGN KEY (address_id) REFERENCES user_addresses(id)
);
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE,
    invoice_number VARCHAR(50) UNIQUE,
    total_amount NUMERIC(10, 2),
    issued_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT fk_invoice_order FOREIGN KEY (order_id) REFERENCES orders(id)
);
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255),
    message TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_type VARCHAR(20),
    -- percent / flat
    discount_value NUMERIC(10, 2),
    min_order_value NUMERIC(10, 2),
    max_discount NUMERIC(10, 2),
    usage_limit INT,
    used_count INT DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE IF NOT EXISTS coupon_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    coupon_id UUID NOT NULL,
    used_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_id, coupon_id),
    CONSTRAINT fk_coupon_usages_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_coupon_usages_coupon FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_shipments_order_id ON shipments(order_id);
CREATE INDEX IF NOT EXISTS idx_shipments_address_id ON shipments(address_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);