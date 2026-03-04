DROP INDEX IF EXISTS idx_password_reset_user_active;
DROP INDEX IF EXISTS idx_password_reset_user_id;
DROP INDEX IF EXISTS idx_single_default_address;
DROP INDEX IF EXISTS idx_user_addresses_user_id;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS ux_users_email;
DROP INDEX IF EXISTS ux_users_username;

DROP TABLE IF EXISTS password_reset;
DROP TABLE IF EXISTS user_addresses;
DROP TABLE IF EXISTS users;