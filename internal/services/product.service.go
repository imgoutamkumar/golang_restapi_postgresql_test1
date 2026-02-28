package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

func GetProductImages(productId string) ([]dto.ProductImageResponse, error) {
	images, err := repository.GetImagesByProductID(productId)

	responseImages := []dto.ProductImageResponse{}

	for _, img := range images {
		responseImages = append(responseImages, dto.ProductImageResponse{
			Id:        img.ID.String(),
			URL:       img.ImageURL,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicID,
		})
	}

	return responseImages, err
}
func GetVariantWithCache(variantId string) (*models.ProductVariant, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product_variant:details:%s", variantId)
	val, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var productVariant models.ProductVariant
		if jsonErr := json.Unmarshal([]byte(val), &productVariant); jsonErr == nil {
			return &productVariant, nil
		}
	}
	if err != nil {
		fmt.Println("Redis error:", err)
	}
	return nil, errors.New("variant not found in cache")
}

func CreateProductService(c *gin.Context, req *dto.CreateProductRequest, userId uuid.UUID) (*models.Product, error) {
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	product := models.Product{
		Name:       req.Name,
		BrandID:    uuid.MustParse(req.BrandID),
		CategoryID: uuid.MustParse(req.CategoryID),
		Status:     models.ProductActive,
		CreatedBy:  userId,
	}
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// variantAttrMap := make(map[uuid.UUID][]uuid.UUID)
	uploadedPublicIDs := []string{}
	var variants []dto.CreateVariantRequest
	if err := json.Unmarshal([]byte(req.VariantsJSON), &variants); err != nil {
		return nil, err
	}

	form, formErr := c.MultipartForm()
	if formErr != nil {
		fmt.Println("MULTIPART ERROR:", formErr)
		return nil, formErr
	}
	for i, v := range variants {

		variant := models.ProductVariant{
			ProductID:       product.ID,
			Sku:             v.Sku,
			Price:           v.Price,
			Stock:           v.Stock,
			DiscountPercent: v.DiscountPercent,
			IsDefault:       v.IsDefault,
			Status:          models.ProductStatus(v.Status),
		}
		fmt.Println("AttributeValueIDs:", v.AttributeValueIDs)
		if err := tx.Create(&variant).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		var attrRows []models.VariantAttribute
		for _, attrId := range v.AttributeValueIDs {
			parsedID, err := uuid.Parse(attrId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			attrRows = append(attrRows, models.VariantAttribute{
				VariantID:        variant.ID,
				AttributeValueID: parsedID,
			})
		}

		if len(attrRows) > 0 {
			if err := tx.Create(&attrRows).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		key := fmt.Sprintf("variant_images_%d", i)
		files := form.File[key] // []*multipart.FileHeader
		if len(files) == 0 {
			tx.Rollback()
			return nil, errors.New("images required for variant")
		}

		if v.PrimaryIndex >= len(files) {
			tx.Rollback()
			return nil, errors.New("invalid primary index")
		}

		folder := fmt.Sprintf("ecommerce/products/%s", product.ID.String())
		var wg sync.WaitGroup
		var mu sync.Mutex
		var uploadErr error
		var errMu sync.Mutex
		modelImages := make([]models.ProductImage, len(files))
		for imgdx, fileHeader := range files {
			wg.Add(1)

			go func(i int, fileHeader *multipart.FileHeader) {
				defer wg.Done()
				// Stop if already failed
				errMu.Lock()
				if uploadErr != nil {
					errMu.Unlock()
					return
				}
				errMu.Unlock()

				uploadFileData, err := utils.UploadFileToCloudinary(fileHeader, folder)
				fmt.Printf("Upload result for file %s: %v\n", fileHeader.Filename, uploadFileData)
				if err != nil {
					errMu.Lock()
					if uploadErr == nil {
						uploadErr = err
					}
					errMu.Unlock()
					return
				}

				mu.Lock()
				uploadedPublicIDs = append(uploadedPublicIDs, uploadFileData.Public_Id)

				modelImages[i] = models.ProductImage{
					VariantID: variant.ID,
					ImageURL:  uploadFileData.ImageUrl,
					IsPrimary: (i == v.PrimaryIndex), // selected primary
					SortOrder: i,                     // maintain order
					PublicID:  uploadFileData.Public_Id,
				}
				mu.Unlock()

			}(imgdx, fileHeader)
		}

		wg.Wait()

		if uploadErr != nil {
			tx.Rollback()
			for _, publicID := range uploadedPublicIDs {
				utils.DeleteFileFromCloudinary(publicID)
			}
			return nil, uploadErr
		}
		if err := tx.Create(&modelImages).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		variant.Images = modelImages
		product.Variants = append(product.Variants, variant)
	}
	// err := repository.CreateProduct(&product, variantAttrMap)
	// err := repository.CreateNewProduct(&product)
	// ✅ FINAL COMMIT
	err := tx.Commit().Error
	if err != nil {
		for _, publicID := range uploadedPublicIDs {
			utils.DeleteFileFromCloudinary(publicID)
		}
		return nil, err
	}
	return &product, nil
}

func GetVariantById(productId string, variantId string) (*dto.ProductVariantResponse, error) {

	// product, err := repository.GetProductByUUID(productID)
	productVariant, err := repository.GetVariantById(productId, variantId)

	if err != nil {
		return nil, err
	}
	// response := utils.MapVariantToResponse(productVariant)
	return productVariant, nil
}
