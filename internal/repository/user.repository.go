package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func Register(user *models.User) (*models.User, error) {
	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUUID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Select("id", "username", "email", "created_at").First(&user, "email = ?", email).Error; err != nil {
		return nil, err // GORM returns gorm.ErrRecordNotFound if no match
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	result := config.DB.
		Select("id", "username", "email", "created_at").
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
