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
			URL:       img.ImageUrl,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicId,
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

func CreateProductService(c *gin.Context, req *dto.CreateProductRequest, userId string) (*models.Product, error) {
	product := models.Product{
		Name:       req.Name,
		BrandID:    req.BrandID,
		CategoryID: req.CategoryID,
		Status:     models.ProductActive,
	}

	// variantAttrMap := make(map[uuid.UUID][]uuid.UUID)
	uploadedURLs := []string{}
	for _, v := range req.Variants {

		if v.PrimaryIndex >= len(v.ImageFiles) {
			return nil, errors.New("invalid primary image index")
		}

		variant := models.ProductVariant{
			// ID:    uuid.New(),
			Sku:             v.Sku,
			Price:           v.Price,
			Stock:           v.Stock,
			DiscountPercent: v.DiscountPercent,
			IsDefault:       v.IsDefault,
			Status:          models.ProductStatus(v.Status),
		}

		for _, attrId := range v.AttributeValueIDs {
			variant.VariantAttributes = append(variant.VariantAttributes, models.VariantAttribute{
				ID:               uuid.New(),
				VariantID:        variant.ID,
				AttributeValueID: attrId,
			})
		}
		folder := fmt.Sprintf("ecommerce/products/%s", product.ID.String())
		var wg sync.WaitGroup
		var mu sync.Mutex
		var uploadErr error
		modelImages := make([]models.ProductImage, len(v.ImageFiles))
		for i, fileHeader := range v.ImageFiles {
			wg.Add(1)

			go func(i int, fileHeader *multipart.FileHeader) {
				defer wg.Done()

				uploadFileData, err := utils.UploadFileToCloudinary(fileHeader, folder)
				if err != nil {
					uploadErr = err
					return
				}

				mu.Lock()
				uploadedURLs = append(uploadedURLs, uploadFileData.ImageUrl)

				modelImages[i] = models.ProductImage{
					ImageURL:  uploadFileData.ImageUrl,
					IsPrimary: (i == v.PrimaryIndex), // selected primary
					SortOrder: i,                     // maintain order
					PublicID:  uploadFileData.Public_Id,
				}
				mu.Unlock()

			}(i, fileHeader)
		}

		wg.Wait()

		if uploadErr != nil {
			for _, url := range uploadedURLs {
				utils.DeleteFile(url)
			}
			utils.ResponseError(c, 500, "Image upload failed", uploadErr)
			return nil, uploadErr
		}

		variant.Images = modelImages
		product.Variants = append(product.Variants, variant)
	}
	// err := repository.CreateProduct(&product, variantAttrMap)
	err := repository.CreateNewProduct(&product)
	if err != nil {
		for _, url := range uploadedURLs {
			utils.DeleteFile(url)
		}
		return nil, err
	}
	return nil, err
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
