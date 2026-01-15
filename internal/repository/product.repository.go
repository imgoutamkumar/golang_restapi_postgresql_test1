package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func CreateProduct(product *models.Product) (*models.Product, error) {
	if err := config.DB.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func GetAllProduct() ([]models.Product, error) {
	var products []models.Product
	result := config.DB.Select("id", "name", "description", "price", "created_at").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func GetProductByUUID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := config.DB.Preload("User").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func UpdateProduct(product *models.Product) error {
	return config.DB.
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(product).
		Error
}
