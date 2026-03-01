package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
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
	userIdStr, _ := c.Get("userId")
	var userId *uuid.UUID
	if userIdStr != nil {
		parsedUserId, err := uuid.Parse(userIdStr.(string))
		if err == nil {
			userId = &parsedUserId
		}
	}
	products, total, err := repository.GetAllProducts(userId, page, limit, search, brand, minPrice, maxPrice, discount)
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
	var req dto.CreateProductRequest
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
	fmt.Printf("User ID from context: %s\n", userId.String())

	if err := c.ShouldBind(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}
	// brandUUID, _ := uuid.Parse(req.BrandID)

	// if err := config.Validate.Struct(req); err != nil {
	// 	log.Printf("%+v\n", err)
	// 	utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
	// 	return
	// }

	// if err := helper.PriceValidate(&req, 0); err != nil {
	// 	utils.ResponseError(c, http.StatusBadRequest, "Price validation failed", err)
	// 	return
	// }
	_, err = services.CreateProductService(c, &req, userId)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}
	// response := utils.MapProductToResponse(createdProduct)

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Product created successfully",
		"data":    nil,
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

func GetProductWithCache(productID uuid.UUID) (*dto.ProductResponse, error) {
	// ctx := context.Background()
	// cacheKey := fmt.Sprintf("product:details:%s", productID.String())

	// val, err := config.RDB.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	var product dto.ProductResponse
	// 	if jsonErr := json.Unmarshal([]byte(val), &product); jsonErr == nil {
	// 		return &product, nil
	// 	}
	// }

	// if err != nil {
	// 	fmt.Println("Redis error:", err)
	// }

	product, err := repository.GetProductByUUID(productID)
	if err != nil {
		return nil, err
	}
	// images, err := services.GetProductImages(product.ID)

	response := dto.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		ShortDesc: product.ShortDesc,
		Currency:  product.Currency,
		CreatedBy: product.CreatedBy,
		CreatedAt: product.CreatedAt,
		Brand: dto.BrandResponse{
			ID:   product.Brand.ID,
			Name: product.Brand.Name,
		},
		Variants: product.Variants,
	}

	// here use goroutine to set cache asynchronouslys

	// go func() {
	// 	data, _ := json.Marshal(response)
	// 	config.RDB.Set(ctx, cacheKey, data, 1*time.Hour)
	// }()

	return &response, nil
}

func ProductImagesReorder(c *gin.Context) {
	var req dto.ReorderProductImagesRequest
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

func GetFilters(c *gin.Context) {
	brands, err := repository.GetBrands()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch brands", nil)
		return
	}
	responses := make([]dto.BrandResponse, 0, len(brands))
	for _, brand := range brands {
		responses = append(responses, dto.BrandResponse{
			ID:   brand.ID.String(),
			Name: brand.Name,
		})

	}
	finalResponses := dto.FiltersResponse{
		Brands: responses,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "data fetched in successfully",
		"data":    finalResponses,
	})
}

func CreateNewBrand(c *gin.Context) {
	var req dto.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := repository.CreateBrand(req.Name)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create brand", nil)
		return
	}
	// response := dto.BrandResponse{
	// 	ID:   createdBrand.ID.String(),
	// 	Name: createdBrand.Name,
	// }
	utils.ResponseSuccess(c, http.StatusOK, "brand created successfully", nil)
}

func CreateNewAttribute(c *gin.Context) {
	var req dto.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := repository.CreateAttribute(req.Name)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create attribute", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "attribute created successfully", nil)
}

func CreateNewAttributeValue(c *gin.Context) {
	attributeId := c.Param("attributeId") // returns string
	id, err := uuid.Parse(attributeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var req dto.CreateAttributeValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = repository.CreateAttributeValue(req, id)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create attribute value", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusCreated, "attribute value created successfully", nil)
}

func GetAllAttributes(c *gin.Context) {

}

func GetAllAttributeValues(c *gin.Context) {

}

func GetAttributeValuesByAttributeId(c *gin.Context) {

}

func GetAttributeValueById(c *gin.Context) {

}

func GetAttributeById(c *gin.Context) {

}

func GetNewArrivals(c *gin.Context) {

}

func GetBestSellers(c *gin.Context) {

}

func GetProductsByCategory(c *gin.Context) {

}

func CreateNewCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if err := config.Validate.Struct(req); err != nil {
		log.Printf("%+v\n", err)
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}
	if err := repository.CreateCategory(req); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create category", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusCreated, "category created successfully", nil)
}

func GetAllCategories(c *gin.Context) {

}

func GetCategoryById(c *gin.Context) {

}

func UpdateCategory(c *gin.Context) {

}

func DeleteCategory(c *gin.Context) {

}
