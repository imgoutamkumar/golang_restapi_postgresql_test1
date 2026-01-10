DROP TRIGGER IF EXISTS cart_updated_at ON cart;
DROP TRIGGER IF EXISTS orders_updated_at ON orders;
DROP TRIGGER IF EXISTS products_updated_at ON products;
DROP TRIGGER IF EXISTS users_updated_at ON users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS cart CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS order_status;
