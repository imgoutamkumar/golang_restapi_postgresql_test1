CREATE TYPE inventory_reason AS ENUM (
    'RESTOCK',
    'ORDER_PLACED',
    'ORDER_CANCELLED',
    'RETURNED',
    'DAMAGED',
    'MANUAL_ADJUSTMENT'
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
    variant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    -- Stock ready to be sold
    available_stock INT NOT NULL DEFAULT 0 CHECK (available_stock >= 0),
    -- Stock currently in someone's checkout process, but not yet paid for
    reserved_stock INT NOT NULL DEFAULT 0 CHECK (reserved_stock >= 0),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Ensure we only have one row per variant per warehouse
    UNIQUE (variant_id, warehouse_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_inventory_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    CONSTRAINT fk_inventory_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
);
-- 3. Inventory Ledgers (The Immutable Audit Trail)
-- An append-only log of EVERY time stock goes up or down. 
-- NEVER UPDATE rows in this table, only INSERT.
CREATE TABLE IF NOT EXISTS inventory_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    -- Positive numbers (+50) for restock/returns. Negative numbers (-1) for sales/damage.
    quantity_change INT NOT NULL,
    -- Why did the stock change? 
    -- e.g., 'RESTOCK', 'ORDER_PLACED', 'RETURNED_TO_STOCK', 'DAMAGED_IN_TRANSIT'
    reason inventory_reason NOT NULL,
    -- Links to an Order ID or a Purchase Order ID for perfect traceability
    reference_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ledger_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    CONSTRAINT fk_ledger_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
);
-- Index for fast lookups when querying a product's stock
CREATE INDEX IF NOT EXISTS idx_inventory_levels_variant ON inventory_levels(variant_id);
-- Index for auditing and calculating historical stock
CREATE INDEX IF NOT EXISTS idx_inventory_ledgers_variant ON inventory_ledgers(variant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_ledgers_created_at ON inventory_ledgers(created_at);
CREATE INDEX IF NOT EXISTS idx_inventory_levels_warehouse ON inventory_levels(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_levels_variant_warehouse ON inventory_levels(variant_id, warehouse_id);
-- Ensure only one primary image per product    