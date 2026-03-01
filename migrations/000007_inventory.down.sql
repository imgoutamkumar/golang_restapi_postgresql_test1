-- =========================
-- DOWN MIGRATION: INVENTORY
-- =========================

-- 1. Drop indexes
DROP INDEX IF EXISTS idx_inventory_levels_variant;
DROP INDEX IF EXISTS idx_inventory_ledgers_variant;
DROP INDEX IF EXISTS idx_inventory_ledgers_created_at;
DROP INDEX IF EXISTS idx_inventory_levels_warehouse;
DROP INDEX IF EXISTS idx_inventory_levels_variant_warehouse;

-- 2. Drop tables (child → parent order)

-- Ledger depends on inventory_levels indirectly via FKs
DROP TABLE IF EXISTS inventory_ledgers;

-- Inventory levels depends on warehouses + product_variants
DROP TABLE IF EXISTS inventory_levels;

-- Warehouses depends on users
DROP TABLE IF EXISTS warehouses;

-- 3. Drop ENUM type (must be last)
DROP TYPE IF EXISTS inventory_reason;