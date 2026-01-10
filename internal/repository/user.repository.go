package repository

import (
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func Register(user *models.User) error {
	return config.DB.Create(user).Error
}
