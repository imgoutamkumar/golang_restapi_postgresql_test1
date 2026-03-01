CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
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
CREATE UNIQUE INDEX ux_users_username_lower ON users (LOWER(username))
WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX ux_users_email_lower ON users (LOWER(email))
WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role_id ON users(role_id);
-- 1. Standard lookup index (Used every time a user opens their "My Addresses" page)
CREATE INDEX IF NOT EXISTS idx_user_addresses_user_id ON user_addresses(user_id)
WHERE deleted_at IS NULL;
-- 2. The FAANG "Single Default" Constraint (Partial Unique Index)
-- This guarantees that a user can ONLY have ONE default address at a time.
-- If you try to set two addresses to is_default=true, PostgreSQL will block it.
CREATE UNIQUE INDEX IF NOT EXISTS idx_single_default_address ON user_addresses(user_id)
WHERE is_default = true
    AND deleted_at IS NULL;
-- Password reset safety
CREATE UNIQUE INDEX idx_password_reset_user_active ON password_reset(user_id)
WHERE locked_at IS NULL
    AND expires_at > now();