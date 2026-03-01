-- function (safe)
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- users
DROP TRIGGER IF EXISTS users_updated_at ON users;
CREATE TRIGGER users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- products
DROP TRIGGER IF EXISTS products_updated_at ON products;
CREATE TRIGGER products_updated_at BEFORE
UPDATE ON products FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- product_images
DROP TRIGGER IF EXISTS product_images_updated_at ON product_images;
CREATE TRIGGER product_images_updated_at BEFORE
UPDATE ON product_images FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- carts
DROP TRIGGER IF EXISTS carts_updated_at ON carts;
CREATE TRIGGER carts_updated_at BEFORE
UPDATE ON carts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- orders
DROP TRIGGER IF EXISTS orders_updated_at ON orders;
CREATE TRIGGER orders_updated_at BEFORE
UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();