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


// Database = 100% empty.
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

date 17-02-26  handler for create new product use this

product := models.Product{
	Name:       req.Name,
	BrandID:    req.BrandID,
	CategoryID: req.CategoryID,
	Status:     models.ProductActive,
}

variantAttrMap := make(map[uuid.UUID][]uuid.UUID)

for _, v := range req.Variants {

	variant := models.ProductVariant{
		ID:    uuid.New(),
		Sku:   v.Sku,
		Price: v.Price,
		Stock: v.Stock,
	}

	// images
	for _, img := range v.Images {
		variant.Images = append(variant.Images, models.ProductImage{
			ImageURL:  img.ImageURL,
			PublicID:  img.PublicID,
			IsPrimary: img.IsPrimary,
		})
	}

	product.Variants = append(product.Variants, variant)

	// map variant -> attribute ids
	variantAttrMap[variant.ID] = v.AttributeValueIDs
}

created, err := repo.CreateProduct(&product, variantAttrMap)