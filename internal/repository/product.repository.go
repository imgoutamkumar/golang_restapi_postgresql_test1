package repository

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

type productRow struct {
	ID               string
	Name             string
	ShortDescription string
	BasePrice        float64
	DiscountPercent  float64
	FinalPrice       float64
	Currency         string
	Stock            int
	CreatedAt        time.Time
	CreatedBy        string
	BrandID          string
	BrandName        string
}

func CreateProduct(product *models.Product) (*models.Product, error) {

	db := config.DB
	if err := db.Create(product).Error; err != nil {
		return nil, err
	}
	// reload with relations needed for response
	if err := db.
		Preload("Brand").
		Preload("ProductImages").
		First(product, "id = ?", product.ID).Error; err != nil {
		return nil, err
	}

	return product, nil

}

func GetAllProducts(
	page string,
	limit string,
	search string,
	brand string,
	minPrice string,
	maxPrice string,
	discount string,
) ([]dto.ProductResponse, int64, error) {

	// var products []dto.ProductResponse
	var total int64

	db := config.DB.Table("products").
		Select(`
			products.id,
			products.name,
			products.short_description,
			products.base_price as base_price,
			products.discount_percent,
			(products.base_price - (products.base_price * products.discount_percent / 100)) as final_price,
			products.currency,
			products.number_of_stock as stock,
			products.created_at,
			brands.id as brand_id,
		brands.name as brand_name
		`).
		Joins("LEFT JOIN brands ON brands.id = products.brand_id")

		// ---------------- SEARCH ----------------
	if search != "" {
		search = strings.ToLower(strings.TrimSpace(search))

		db = db.Where(`
		LOWER(products.name) ILIKE ? OR 
		LOWER(products.short_description) ILIKE ?
	`, "%"+search+"%", "%"+search+"%")
	}

	//if fronend send brand_id then no need to add join here
	if brand != "" {
		rawBrands := strings.Split(brand, ",")
		var brands []string

		for _, b := range rawBrands {
			brands = append(brands, strings.TrimSpace(b))
		}

		db = db.Where("LOWER(brands.name) IN ?", brands)
	}

	if minPrice != "" && maxPrice != "" {
		min, _ := strconv.Atoi(minPrice)
		max, _ := strconv.Atoi(maxPrice)
		db = db.Where("base_price BETWEEN ? AND ?", min, max)
	}

	if discount != "" {
		d, _ := strconv.Atoi(discount)
		db = db.Where("discount_percent >= ?", d)
	}

	db.Count(&total)

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	if p <= 0 {
		p = 1
	}
	if l <= 0 {
		l = 12
	}

	offset := (p - 1) * l

	var rows []productRow
	// ---------- fetch images in ONE query ----------
	err := db.
		Limit(l).
		Offset(offset).
		Order("created_at DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	// ------------------ GET ALL PRODUCT IDS ------------------
	var productIDs []string
	for _, r := range rows {
		productIDs = append(productIDs, r.ID)
	}

	// ------------------ FETCH ALL IMAGES IN ONE QUERY ------------------
	var imgRows []struct {
		Id        string
		ProductID string
		ImageURL  string
		IsPrimary bool
		PublicID  string
	}

	config.DB.Table("product_images").
		Select("id, product_id, image_url, is_primary, public_id").
		Where("product_id IN ?", productIDs).
		Order("sort_order ASC").
		Scan(&imgRows)

		// ------------------ BUILD IMAGE MAP ------------------
	imageProductMap := make(map[string][]dto.ProductImageResponse)

	for _, img := range imgRows {
		imageProductMap[img.ProductID] = append(
			imageProductMap[img.ProductID],
			dto.ProductImageResponse{
				Id:        img.Id,
				URL:       img.ImageURL,
				IsPrimary: img.IsPrimary,
				PublicId:  img.PublicID,
			},
		)
	}

	// ---------- build DTO ----------
	responses := make([]dto.ProductResponse, 0, len(rows))

	for _, r := range rows {
		responses = append(responses, dto.ProductResponse{
			ID:              r.ID,
			Name:            r.Name,
			ShortDesc:       r.ShortDescription,
			BasePrice:       r.BasePrice,
			DiscountPercent: r.DiscountPercent,
			FinalPrice:      r.FinalPrice,
			Currency:        r.Currency,
			Stock:           r.Stock,
			CreatedAt:       r.CreatedAt,
			Brand: dto.BrandResponse{
				ID:   r.BrandID,
				Name: r.BrandName,
			},
			Images: imageProductMap[r.ID], // attach images
		})
	}

	return responses, total, nil
}

func GetProductByUUID(id uuid.UUID) (*productRow, error) {
	var row productRow
	db := config.DB
	err := db.Table("products").
		Select(`
			products.id, 
			products.name,
			products.short_description,
			products.base_price as base_price,
			products.discount_percent,
			(products.base_price - (products.base_price * products.discount_percent / 100)) as final_price,
			products.currency,
			products.number_of_stock as stock,
			products.created_by,
			products.created_at,
			brands.id as brand_id,
		brands.name as brand_name
		`).
		Joins("LEFT JOIN brands ON brands.id = products.brand_id").
		Where("products.id = ?", id).Scan(&row).Error

	// if err != nil {
	// 	return nil, err
	// }

	// var imgRows []struct {
	// 	Id        string
	// 	ProductID string
	// 	ImageURL  string
	// 	IsPrimary bool
	// 	PublicID  string
	// }

	// db.Table("product_images").
	// 	Select("id, product_id, image_url, is_primary, public_id").
	// 	Where("product_id = ?", id).
	// 	Order("sort_order ASC").
	// 	Scan(&imgRows)

	// imageProductMap := []dto.ProductImageResponse{}

	// for _, img := range imgRows {
	// 	imageProductMap = append(
	// 		imageProductMap,
	// 		dto.ProductImageResponse{
	// 			Id:        img.Id,
	// 			URL:       img.ImageURL,
	// 			IsPrimary: img.IsPrimary,
	// 			PublicId:  img.PublicID,
	// 		},
	// 	)
	// }

	// response := dto.ProductResponse{
	// 	ID:              row.ID,
	// 	Name:            row.Name,
	// 	ShortDesc:       row.ShortDescription,
	// 	BasePrice:       row.BasePrice,
	// 	DiscountPercent: row.DiscountPercent,
	// 	FinalPrice:      row.FinalPrice,
	// 	Currency:        row.Currency,
	// 	Stock:           row.Stock,
	// 	CreatedAt:       row.CreatedAt,
	// 	Brand: dto.BrandResponse{
	// 		ID:   row.BrandID,
	// 		Name: row.BrandName,
	// 	},
	// 	Images: imageProductMap,
	// }

	return &row, err
}

func UpdateProduct(product *models.Product) error {
	return config.DB.
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(product).
		Error
}

// for transactional purposes
func UpdateStock(db *gorm.DB, productID uuid.UUID, qty int) error {
	return db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", qty)).
		Error
}

func CreateNewProductOptimalApproach(product *models.Product) (*models.Product, error) {
	db := config.DB

	tx := db.Begin()

	// create product
	if err := tx.Create(product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// reload only required fields
	var result models.Product
	if err := db.
		Preload("Brand", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name") // select only needed fields
		}).
		Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "image_url", "is_primary", "sort_order")
		}).
		First(&result, "id = ?", product.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func GetImagesByProductID(productID string) ([]models.ProductImages, error) {

	var images []models.ProductImages

	err := config.DB.Table("product_images").
		Select("id, image_url, is_primary, public_id").
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Scan(&images).Error

	if err != nil {
		return nil, err
	}
	return images, nil
}

func ReorderProductImages(req helper.ReorderProductImagesRequest) error {
	tx := config.DB.Begin()

	for _, img := range req.Images {
		err := tx.Exec(`
			UPDATE product_images
			SET sort_order = ?
			WHERE id = ? AND product_id = ?
		`, img.SortOrder, img.ID, req.ProductID).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
