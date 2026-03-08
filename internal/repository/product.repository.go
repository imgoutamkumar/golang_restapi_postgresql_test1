package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

func CreateProduct(product *models.Product, variantAttrMap map[uuid.UUID][]uuid.UUID) error {
	db := config.DB

	err := db.Transaction(func(tx *gorm.DB) error {

		// create product + variants + images
		if err := tx.
			Omit("Variants.AttributeValues.*").
			Create(product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		// insert pivot table (variant_attributes)
		type Pivot struct {
			VariantID        uuid.UUID `gorm:"type:uuid"`
			AttributeValueID uuid.UUID `gorm:"type:uuid"`
		}

		var pivots []Pivot

		for _, variant := range product.Variants {

			attrIDs := variantAttrMap[variant.ID]
			if len(attrIDs) == 0 {
				continue
			}

			for _, attrID := range attrIDs {
				pivots = append(pivots, Pivot{
					VariantID:        variant.ID,
					AttributeValueID: attrID,
				})
			}
		}

		if len(pivots) > 0 {
			if err := tx.Table("variant_attributes").Create(&pivots).Error; err != nil {
				return fmt.Errorf("failed to insert variant attributes: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// reload full product

	/* var created models.Product

	if err := db.
		Preload("Brand").
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Images").
		Preload("Variants.AttributeValues").
		First(&created, "id = ?", product.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload product: %w", err)
	}

	return &created, nil
	*/
	return err
}

func CreateNewProduct(product *models.Product) error {
	db := config.DB

	err := db.Transaction(func(tx *gorm.DB) error {

		// create product + variants + images
		if err := tx.
			Omit("Variants.AttributeValues.*").
			Create(product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return err
}

func GetAllProducts(
	userID *uuid.UUID,
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

	// ---------------- BASE QUERY ----------------
	query := config.DB.Model(&models.Product{}).
		Joins("LEFT JOIN brands ON brands.id = products.brand_id").
		Where("products.deleted_at IS NULL")

	// ---------------- SEARCH ----------------
	if search != "" {
		search = strings.ToLower(strings.TrimSpace(search))
		query = query.Where(`
			LOWER(products.name) ILIKE ? OR 
			LOWER(products.description) ILIKE ?
		`, "%"+search+"%", "%"+search+"%")
	}

	//if fronend send brand_id then no need to add join here
	if brand != "" {
		rawBrands := strings.Split(brand, ",")
		var brands []string

		for _, b := range rawBrands {
			brands = append(brands, strings.TrimSpace(b))
		}

		query = query.Where("LOWER(brands.name) IN ?", brands)
	}

	// ---------------- COUNT ----------------
	query.Count(&total)

	// ---------------- PAGINATION ----------------
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	if p <= 0 {
		p = 1
	}
	if l <= 0 {
		l = 10
	}

	offset := (p - 1) * l

	// ---------------- FETCH PRODUCTS ----------------
	var products []models.Product
	err := query.
		Preload("Brand").
		Order("products.created_at DESC").
		Limit(l).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return []dto.ProductResponse{}, total, nil
	}

	// ---------------- COLLECT PRODUCT IDS ----------------
	productIDs := make([]uuid.UUID, 0)
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}
	fmt.Println("All Product ids", productIDs)
	// ---------------- FETCH VARIANTS ----------------
	var variants []models.ProductVariant
	err = config.DB.
		Where("product_id IN (?) AND is_default = ?", productIDs, true).
		Preload("Images", "is_primary = ?", true).
		Find(&variants).Error
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("Fetched variants:", len(variants))
	variantMap := make(map[uuid.UUID][]models.ProductVariant)

	for _, v := range variants {
		variantMap[v.ProductID] = append(variantMap[v.ProductID], v)
	}
	fmt.Println("Fetched variantMap:", variantMap)

	// ---------------- FETCH WISHLISTED VARIANT IDS IF USER LOGGED IN ----------------
	wishlistMap := make(map[uuid.UUID]bool)

	if userID != nil {

		type WishlistProduct struct {
			ProductID uuid.UUID
		}

		var wishlistProducts []WishlistProduct

		err = config.DB.
			Table("wishlist_items wi").
			Select("DISTINCT pv.product_id").
			Joins("JOIN product_variants pv ON pv.id = wi.variant_id").
			Joins("JOIN wishlists w ON w.id = wi.wishlist_id").
			Where("w.user_id = ?", *userID).
			Scan(&wishlistProducts).Error

		if err != nil {
			return nil, 0, err
		}

		for _, item := range wishlistProducts {
			wishlistMap[item.ProductID] = true
		}
	}
	// ---------------- BUILD RESPONSE ----------------
	var responses []dto.ProductResponse

	for _, p := range products {

		productResp := dto.ProductResponse{
			ID:        p.ID.String(),
			Name:      p.Name,
			ShortDesc: p.ShortDescription,
			Currency:  p.Currency,
			CreatedAt: p.CreatedAt,
			Brand: dto.BrandResponse{
				ID:   p.Brand.ID.String(),
				Name: p.Brand.Name,
			},
			IsWishlisted: wishlistMap[p.ID],
		}

		var variantResponses []dto.ProductVariantResponse
		// wishlistedVariantIds, err := services.GetWishlistedVariantIds()
		if err != nil {
			fmt.Println("Error fetching wishlisted variant IDs:", err)
		}
		for _, v := range variants {
			fmt.Println("Processing variant:", v)
			var variantImages []dto.ProductImageResponse
			for _, img := range v.Images {
				variantImages = append(variantImages, dto.ProductImageResponse{
					Id:        img.ID.String(),
					URL:       img.ImageURL,
					IsPrimary: img.IsPrimary,
					PublicId:  img.PublicID,
				})
			}
			variantResponses = append(variantResponses,
				dto.ProductVariantResponse{
					Sku:             v.Sku,
					Price:           v.Price,
					DiscountPercent: v.DiscountPercent,
					FinalPrice:      v.Price - (v.Price * v.DiscountPercent / 100),
					Stock:           v.Stock,
					Images:          variantImages,
				})
		}

		productResp.Variants = variantResponses
		responses = append(responses, productResp)
	}

	return responses, total, nil
}

func GetAllProductsForAdmin(
	userID *uuid.UUID,
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

	// ---------------- BASE QUERY ----------------
	query := config.DB.Model(&models.Product{}).
		Joins("LEFT JOIN brands ON brands.id = products.brand_id").
		Where("products.deleted_at IS NULL")

	// ---------------- SEARCH ----------------
	if search != "" {
		search = strings.ToLower(strings.TrimSpace(search))
		query = query.Where(`
			LOWER(products.name) ILIKE ? OR 
			LOWER(products.description) ILIKE ?
		`, "%"+search+"%", "%"+search+"%")
	}

	//if fronend send brand_id then no need to add join here
	if brand != "" {
		rawBrands := strings.Split(brand, ",")
		var brands []string

		for _, b := range rawBrands {
			brands = append(brands, strings.TrimSpace(b))
		}

		query = query.Where("LOWER(brands.name) IN ?", brands)
	}

	// ---------------- COUNT ----------------
	query.Count(&total)

	// ---------------- PAGINATION ----------------
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	if p <= 0 {
		p = 1
	}
	if l <= 0 {
		l = 10
	}

	offset := (p - 1) * l

	// ---------------- FETCH PRODUCTS ----------------
	var products []models.Product
	err := query.
		Preload("Brand").
		Order("products.created_at DESC").
		Limit(l).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return []dto.ProductResponse{}, total, nil
	}

	// ---------------- COLLECT PRODUCT IDS ----------------
	productIDs := make([]uuid.UUID, 0)
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}
	fmt.Println("All Product ids", productIDs)
	// ---------------- FETCH VARIANTS ----------------
	var variants []models.ProductVariant
	err = config.DB.
		Where("product_id IN (?)", productIDs).
		Preload("Images", "is_primary = ?", true).
		Find(&variants).Error
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("Fetched variants:", len(variants))
	variantMap := make(map[uuid.UUID][]models.ProductVariant)

	for _, v := range variants {
		variantMap[v.ProductID] = append(variantMap[v.ProductID], v)
	}
	fmt.Println("Fetched variantMap:", variantMap)

	// ---------------- FETCH WISHLISTED VARIANT IDS IF USER LOGGED IN ----------------
	wishlistMap := make(map[uuid.UUID]bool)

	if userID != nil {

		type WishlistProduct struct {
			ProductID uuid.UUID
		}

		var wishlistProducts []WishlistProduct

		err = config.DB.
			Table("wishlist_items wi").
			Select("DISTINCT pv.product_id").
			Joins("JOIN product_variants pv ON pv.id = wi.variant_id").
			Joins("JOIN wishlists w ON w.id = wi.wishlist_id").
			Where("w.user_id = ?", *userID).
			Scan(&wishlistProducts).Error

		if err != nil {
			return nil, 0, err
		}

		for _, item := range wishlistProducts {
			wishlistMap[item.ProductID] = true
		}
	}
	// ---------------- BUILD RESPONSE ----------------
	var responses []dto.ProductResponse

	for _, p := range products {

		productResp := dto.ProductResponse{
			ID:        p.ID.String(),
			Name:      p.Name,
			ShortDesc: p.ShortDescription,
			Currency:  p.Currency,
			CreatedAt: p.CreatedAt,
			Brand: dto.BrandResponse{
				ID:   p.Brand.ID.String(),
				Name: p.Brand.Name,
			},
			IsWishlisted: wishlistMap[p.ID],
		}

		var variantResponses []dto.ProductVariantResponse
		// wishlistedVariantIds, err := services.GetWishlistedVariantIds()
		if err != nil {
			fmt.Println("Error fetching wishlisted variant IDs:", err)
		}
		for _, v := range variants {
			fmt.Println("Processing variant:", v)
			var variantImages []dto.ProductImageResponse
			for _, img := range v.Images {
				variantImages = append(variantImages, dto.ProductImageResponse{
					Id:        img.ID.String(),
					URL:       img.ImageURL,
					IsPrimary: img.IsPrimary,
					PublicId:  img.PublicID,
				})
			}
			variantResponses = append(variantResponses,
				dto.ProductVariantResponse{
					Sku:             v.Sku,
					Price:           v.Price,
					DiscountPercent: v.DiscountPercent,
					FinalPrice:      v.Price - (v.Price * v.DiscountPercent / 100),
					Stock:           v.Stock,
					Images:          variantImages,
				})
		}

		productResp.Variants = variantResponses
		responses = append(responses, productResp)
	}

	return responses, total, nil
}

func GetProductByUUID(id uuid.UUID) (*dto.ProductResponse, error) {
	var product models.Product
	db := config.DB

	// 1️⃣ Fetch Product + Brand (ONLY)
	if err := db.
		Preload("Brand").
		First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var variants []models.ProductVariant

	if err := db.
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("sort_order ASC")
		}).
		Preload("VariantAttributes.AttributeValue.AttributeType").
		Where("product_id = ?", id).
		Find(&variants).Error; err != nil {
		return nil, err
	}

	response := dto.ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		ShortDesc: product.ShortDescription,
		BasePrice: product.BasePrice,
		Brand: dto.BrandResponse{
			ID:   product.Brand.ID.String(),
			Name: product.Brand.Name,
		},
	}

	for _, v := range variants {
		var variantImages []dto.ProductImageResponse
		for _, img := range v.Images {
			variantImages = append(variantImages, dto.ProductImageResponse{
				Id:        img.ID.String(),
				URL:       img.ImageURL,
				IsPrimary: img.IsPrimary,
				PublicId:  img.PublicID,
			})
		}
		fmt.Println("variants Attributes", v.VariantAttributes)
		// build attribute groups
		attrGroupMap := make(map[string][]dto.AttributeValueResponse)
		for _, va := range v.VariantAttributes {
			attrGroupName := va.AttributeValue.AttributeType.Name
			attrGroupMap[attrGroupName] = append(attrGroupMap[attrGroupName], dto.AttributeValueResponse{
				ID:    va.AttributeValue.ID.String(),
				Value: va.AttributeValue.Value,
			})
		}
		var attributeGroups []dto.AttributeGroup
		for name, values := range attrGroupMap {
			attributeGroups = append(attributeGroups, dto.AttributeGroup{
				Name:   name,
				Values: values,
			})
		}

		response.Variants = append(response.Variants, dto.ProductVariantResponse{
			Sku:            v.Sku,
			Price:          v.Price,
			Stock:          v.Stock,
			Images:         variantImages,
			AttributeGroup: attributeGroups,
		})

	}
	fmt.Println("response", response)
	return &response, nil
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

func GetImagesByProductID(productID string) ([]models.ProductImage, error) {

	var images []models.ProductImage

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

func ReorderProductImages(req dto.ReorderProductImagesRequest) error {
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

func GetBrands() ([]models.Brand, error) {
	var brands []models.Brand
	err := config.DB.Model(&models.Brand{}).Find(&brands).Error
	return brands, err
}

func CreateBrand(name string) error {
	brand := models.Brand{
		Name: name,
	}
	err := config.DB.Create(&brand).Error
	return err
}

func CreateAttribute(name string) error {
	attribute := models.AttributeType{
		Name: name,
	}
	err := config.DB.Create(&attribute).Error
	return err
}

func CreateAttributeValue(req dto.CreateAttributeValueRequest, typeId uuid.UUID) error {
	attributeValue := models.AttributeValue{
		Value:           req.Value,
		AttributeTypeID: typeId,
		MetaInfo:        req.MetaInfo,
	}
	err := config.DB.Create(&attributeValue).Error
	return err
}

func GetAttributeTypes() ([]models.AttributeType, error) {
	var types []models.AttributeType

	err := config.DB.
		Model(&models.AttributeType{}).
		Find(&types).Error

	if err != nil {
		return nil, err
	}

	return types, nil
}

func GetAttributeValuesByAttributeID(attributeID uuid.UUID) ([]dto.AttributeValueResponse, error) {
	var values []dto.AttributeValueResponse

	err := config.DB.
		Model(&models.AttributeValue{}).
		Select("id, value, meta_info").
		Where("attribute_type_id = ?", attributeID).
		Find(&values).Error

	if err != nil {
		return nil, err
	}

	return values, nil
}

func GetVariantById(productId string, variantId string) (*dto.ProductVariantResponse, error) {
	var variant models.ProductVariant
	err := config.DB.Where("product_id = ? AND id = ?", productId, variantId).
		Preload("Images").
		Preload("VariantAttributes.AttributeValue.AttributeType").
		First(&variant).Error
	if err != nil {
		return nil, err
	}

	var productvariant dto.ProductVariantResponse = dto.ProductVariantResponse{
		Sku:             variant.Sku,
		Price:           variant.Price,
		DiscountPercent: variant.DiscountPercent,
		FinalPrice:      variant.Price - (variant.Price * variant.DiscountPercent / 100),
		Stock:           variant.Stock,
	}

	var variantImages []dto.ProductImageResponse
	for _, img := range variant.Images {
		variantImages = append(variantImages, dto.ProductImageResponse{
			Id:        img.ID.String(),
			URL:       img.ImageURL,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicID,
		})
	}

	var attributeGroups []dto.AttributeGroup
	for _, va := range variant.VariantAttributes {
		attributeGroups = append(attributeGroups, dto.AttributeGroup{
			Name: va.AttributeValue.AttributeType.Name,
			Values: []dto.AttributeValueResponse{
				{
					ID:    va.AttributeValue.ID.String(),
					Value: va.AttributeValue.Value,
				},
			},
		})
	}

	productvariant.Images = variantImages
	productvariant.AttributeGroup = attributeGroups
	return &productvariant, nil
}

func GetVariantByVariantId(variantId string) (*dto.ProductVariantResponse, error) {
	var variant models.ProductVariant
	err := config.DB.Where("id = ?", variantId).
		Preload("Images").
		Preload("VariantAttributes.AttributeValue.AttributeType").
		First(&variant).Error
	if err != nil {
		return nil, err
	}

	var productvariant dto.ProductVariantResponse = dto.ProductVariantResponse{
		Sku:             variant.Sku,
		Price:           variant.Price,
		DiscountPercent: variant.DiscountPercent,
		FinalPrice:      variant.Price - (variant.Price * variant.DiscountPercent / 100),
		Stock:           variant.Stock,
	}

	var variantImages []dto.ProductImageResponse
	for _, img := range variant.Images {
		variantImages = append(variantImages, dto.ProductImageResponse{
			Id:        img.ID.String(),
			URL:       img.ImageURL,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicID,
		})
	}

	var attributeGroups []dto.AttributeGroup
	for _, va := range variant.VariantAttributes {
		attributeGroups = append(attributeGroups, dto.AttributeGroup{
			Name: va.AttributeValue.AttributeType.Name,
			Values: []dto.AttributeValueResponse{
				{
					ID:    va.AttributeValue.ID.String(),
					Value: va.AttributeValue.Value,
				},
			},
		})
	}

	productvariant.Images = variantImages
	productvariant.AttributeGroup = attributeGroups
	return &productvariant, nil
}

func GetNewArrivals() ([]dto.ProductResponse, error) {
	var products []models.Product
	err := config.DB.
		Preload("Brand").
		Order("created_at DESC").
		Limit(10).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, dto.ProductResponse{
			ID:        product.ID.String(),
			Name:      product.Name,
			ShortDesc: product.ShortDescription,
			Currency:  product.Currency,
			CreatedBy: product.CreatedBy.String(),
			CreatedAt: product.CreatedAt,
			Brand: dto.BrandResponse{
				ID:   product.Brand.ID.String(),
				Name: product.Brand.Name,
			},
		})

		variants := []models.ProductVariant{}
		err = config.DB.
			Where("product_id = ? AND created_at > ? AND status = ?", product.ID, time.Now().AddDate(0, 0, -30), "active").
			Preload("Images", "is_primary = ?", true).
			Order("created_at DESC").Limit(10).
			Find(&variants).Error
		if err != nil {
			return nil, err
		}
	}

	return responses, nil
}

func CreateCategory(req dto.CreateCategoryRequest) error {
	if req.ParentID != nil {
		var count int64
		config.DB.Model(&models.Category{}).
			Where("id = ?", *req.ParentID).
			Count(&count)

		if count == 0 {
			return errors.New("parent category not found")
		}
	}
	category := models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}
	err := config.DB.Create(&category).Error
	return err
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := config.DB.Find(&categories).Error
	return categories, err
}

func GetProductIdsFromVariantIds(variantIds []uuid.UUID) ([]uuid.UUID, error) {
	var productIds []uuid.UUID
	err := config.DB.Model(&models.ProductVariant{}).
		Where("id IN ?", variantIds).
		Pluck("DISTINCT product_id", &productIds).Error
	if err != nil {
		return nil, err
	}
	return productIds, nil
}
