#used dependencies

go get -u github.com/gin-gonic/gin                             // framework
go get github.com/go-playground/validator/v10                  // validation
go get github.com/shopspring/decimal                           // decimal


# not used but best choice:

go get github.com/gin-contrib/cors                             //  for cors
go get github.com/gorilla/websocket                            //  for web socket



//insert at start

INSERT INTO roles (name) VALUES ('user');
INSERT INTO roles (name) VALUES ('admin');
INSERT INTO roles (name) VALUES ('seller');


next task
define audience schema 

men
women
kids
unisex

define category schema 

clothing
footwear
electronics
beauty
accessories

define subcategory schema 

tshirt
shirt
jeans
kurta
saree
sneakers
heels


brands
------
id
name (Nike, Puma, Roadster)
logo_url

categories
----------
id
name
type

subcategories
-------------
id
name
category_id (FK categories)

product_variants
----------------
id
product_id
size (S,M,L,XL)
color (red, black)
sku
stock
price