package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/services"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

// use goroutines and channels to handle multiple tasks concurrently
// Scenario: "Load Product Page"
// When a user views an iPhone page, you need:
// Product Details (Name, Price).
// Reviews (4.5 Stars).
// "People also bought" recommendations.
// Estimated Delivery Date.

func GetAllProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "12")
	search := c.Query("searchTerm")
	brand := c.Query("brand") // Nike,Puma
	minPrice := c.Query("minPrice")
	maxPrice := c.Query("maxPrice")
	discount := c.Query("discount")
	products, total, err := repository.GetAllProducts(page, limit, search, brand, minPrice, maxPrice, discount)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products, "message": "Data Fetched successfully", "status": true, "total": total})
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// product, err := repository.GetProductByUUID(productID)
	product, err := GetProductWithCache(productID)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	// response := utils.MapProductToResponse(product)

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product fetched successfully",
		"data":    product,
	})
}

func CreateNewProduct(c *gin.Context) {
	var req helper.CreateProductRequest
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid userid": err})
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid form data", err)
		return
	}

	brandUUID, _ := uuid.Parse(req.BrandID)

	if err := config.Validate.Struct(req); err != nil {
		log.Printf("%+v\n", err)
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := helper.CustomValidate(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed.", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var uploadErr error
	// var modelImages []models.ProductImages
	modelImages := make([]models.ProductImages, len(req.ImageFiles))
	uploadedURLs := []string{}
	for i, fileHeader := range req.ImageFiles {
		wg.Add(1)

		go func(i int, fileHeader *multipart.FileHeader) {
			defer wg.Done()

			uploadFileData, err := utils.UploadFileToCloudinary(fileHeader)
			if err != nil {
				uploadErr = err
				return
			}

			mu.Lock()
			uploadedURLs = append(uploadedURLs, uploadFileData.ImageUrl)

			// modelImages = append(modelImages, models.ProductImages{
			// 	ImageUrl:  uploadFileData.ImageUrl,
			// 	IsPrimary: (i == req.PrimaryIndex),
			// 	SortOrder: i,
			// 	PublicId:  uploadFileData.Public_Id,
			// })

			// assign by index (NOT append)
			modelImages[i] = models.ProductImages{
				ImageUrl:  uploadFileData.ImageUrl,
				IsPrimary: (i == req.PrimaryIndex), // selected primary
				SortOrder: i,                       // maintain order
				PublicId:  uploadFileData.Public_Id,
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
		return
	}

	product := models.Product{
		Name:             req.Name,
		ShortDescription: req.Description,
		BasePrice:        req.BasePrice,
		CreatedBy:        userId,
		DiscountPercent:  req.DiscountPercent,
		NumberOfStock:    req.NumberOfStock,
		BrandID:          brandUUID,
		// Category:         req.Category,
		Currency:       req.Currency,
		Status:         req.Status,
		IsReturnable:   req.IsReturnable,
		IsCODAvailable: req.IsCODAvailable,
		Description:    req.Description,
		ProductImages:  modelImages,
	}
	createdProduct, err := repository.CreateProduct(&product)
	if err != nil {
		for _, url := range uploadedURLs {
			utils.DeleteFile(url)
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Error While creating product", err)
		return
	}
	response := utils.MapProductToResponse(createdProduct)

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product created successfully",
		"data":    response,
	})
}

func UpdateProduct(c *gin.Context) {
	var ctx = context.Background()
	var product models.Product
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	existingProduct, err := repository.GetProductByUUID(productID)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if userId.String() != existingProduct.CreatedBy {
		utils.ResponseError(c, http.StatusInternalServerError, "Not Authorized", err)
		return
	}

	if err := repository.UpdateProduct(&product); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Update failed", err)
		return
	}
	// 2. Invalidate Cache (Delete the specific key)
	cacheKey := fmt.Sprintf("product:details:%s", product.ID.String())
	if err := config.RDB.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Println("Failed to clear cache:", err)
	}
	utils.ResponseSuccess(c, http.StatusOK, "product updated successfully", product)
}

func DeleteProduct(c *gin.Context) {

}

// some helper function to get product with caching
func GetProductWithCache(productID uuid.UUID) (*dto.ProductResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:details:%s", productID.String())

	val, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache Hit! Unmarshal JSON to Struct
		var product dto.ProductResponse
		if jsonErr := json.Unmarshal([]byte(val), &product); jsonErr == nil {
			return &product, nil
		}
	}

	if err != nil {
		fmt.Println("Redis error:", err)
	}

	product, err := repository.GetProductByUUID(productID)
	if err != nil {
		return nil, err
	}
	images, err := services.GetProductImages(product.ID)

	response := dto.ProductResponse{
		ID:              product.ID,
		Name:            product.Name,
		ShortDesc:       product.ShortDescription,
		BasePrice:       product.BasePrice,
		DiscountPercent: product.DiscountPercent,
		FinalPrice:      product.FinalPrice,
		Currency:        product.Currency,
		Stock:           product.Stock,
		CreatedBy:       product.CreatedBy,
		CreatedAt:       product.CreatedAt,
		Brand: dto.BrandResponse{
			ID:   product.BrandID,
			Name: product.BrandName,
		},
		Images: images,
	}

	// here use goroutine to set cache asynchronouslys
	go func() {
		data, _ := json.Marshal(response)
		config.RDB.Set(ctx, cacheKey, data, 1*time.Hour)
	}()

	return &response, nil
}

func ProductImagesReorder(c *gin.Context) {
	var req helper.ReorderProductImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.ReorderProductImages(req)
	if err != nil {
		utils.ResponseSuccess(c, http.StatusInternalServerError, "product images failed to reorder", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "product images reorder successfully", nil)
}
