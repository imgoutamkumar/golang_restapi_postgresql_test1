-- Down migration: rollback roles & permissions

-- 1. Drop indexes (optional, but clean)
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;

-- 2. Drop junction table first (has FKs)
DROP TABLE IF EXISTS role_permissions;

-- 3. Drop permissions table
DROP TABLE IF EXISTS permissions;

-- 4. Drop roles table
DROP TABLE IF EXISTS roles;