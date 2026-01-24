package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

func GetAllProducts(c *gin.Context) {
	products, err := repository.GetAllProduct()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", products)
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	product, err := repository.GetProductByUUID(productID)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", product)
}

func CreateNewProduct(c *gin.Context) {
	var req helper.CreateProductRequest
	// for now i can make createrId static, later we get from auth middleware
	// Read JSON body

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validation check using Validate package
	if err := config.Validate.Struct(req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Custom validation
	if err := helper.CustomValidate(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	product := models.Product{
		Name:             req.Name,
		ShortDescription: req.Description,
		BasePrice:        req.BasePrice,
		// for now static, later get from auth middleware
		CreatedBy: uuid.MustParse("f064c2fc-3523-4e1b-b166-e2cc57064fd8"),
	}
	createdProduct, err := repository.CreateProduct(&product)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "product created successfully", createdProduct)
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
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

	if existingProduct.CreatedBy != uuid.MustParse("f064c2fc-3523-4e1b-b166-e2cc57064fd8") {
		utils.ResponseError(c, http.StatusBadRequest, "You are not authorized", err)
		return
	}

	if err := repository.UpdateProduct(&product); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Update failed", err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "product updated successfully", product)
}

func DeleteProduct(c *gin.Context) {

}
