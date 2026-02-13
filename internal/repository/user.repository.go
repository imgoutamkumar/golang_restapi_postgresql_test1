package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func Login() {}

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
	if err := config.DB.Select("id", "username", "email", "password", "created_at").First(&user, "email = ?", email).Error; err != nil {
		return nil, err // GORM returns gorm.ErrRecordNotFound if no match
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	result := config.DB.
		Select("id", "username", "email", "created_at").
		Find(&users)

	return users, result.Error
}

func FilterAndSearchUsers(params helper.UserFilterParams) (*[]models.User, int64, error) {
	var users []models.User
	var total int64
	offset := (params.Page - 1) * params.Limit
	query := config.DB.Model(&models.User{})
	if params.FullName != "" {
		query = query.Where("fullname ILIKE ?", "%"+params.FullName+"%")
	}

	// 3. Filter by Product Name (Joined via Orders -> OrderItems -> Products)
	if params.ProductName != "" {
		// Path: Users -> Orders -> OrderItems -> Products
		query = query.Joins("JOIN orders ON orders.user_id = users.id").
			Joins("JOIN order_items ON order_items.order_id = orders.id").
			Joins("JOIN products ON order_items.product_id = products.id").
			Where("products.product_name ILIKE ?", "%"+params.ProductName+"%").
			Group("users.id")
	}

	if err := config.DB.Table("(?) as sub", query).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(params.Limit).
		Offset(offset).
		Preload("Orders.OrderItems.Product").
		Find(&users).Error

	return &users, total, err
}

func CreatePasswordReset(reset_password *models.PasswordReset) error {
	reset := &models.PasswordReset{
		UserID:    reset_password.UserID,
		OTPHash:   reset_password.OTPHash,
		ExpiresAt: reset_password.ExpiresAt,
	}
	err := config.DB.Create(reset).Error
	return err
}

func GetPasswordResetByUserID(userID string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := config.DB.Where("user_id = ?", userID).First(&reset).Error
	return &reset, err
}

func UpdatePasswordReset(reset *models.PasswordReset) error {
	return config.DB.Save(reset).Error
}

func DeletePasswordReset(id uint) error {
	return config.DB.Delete(&models.PasswordReset{}, id).Error
}
