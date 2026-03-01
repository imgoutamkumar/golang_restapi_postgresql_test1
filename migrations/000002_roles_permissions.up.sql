CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
-- Insert the default role so it exists
INSERT INTO roles (name)
VALUES ('user') ON CONFLICT DO NOTHING;
-- permission table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
-------------------------------------------------------------------------------
-- Example: Insert a permission and assign it to the 'user' role
-- INSERT INTO permissions (name) VALUES
-- ('order:read')
-- ON CONFLICT DO NOTHING;
-------------------------------------------------------------------------------
-- Then assign this permission to the 'user' role
-- (Uncomment and run after inserting the permission)
-- INSERT INTO role_permissions (role_id, permission_id)
-- SELECT r.id, p.id
-- FROM roles r, permissions p
-- WHERE r.name = 'user'
-- AND p.name IN ('order:read')
-- ON CONFLICT DO NOTHING;
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);