package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar,omitempty"` // optional
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		//Avatar:   user.AvatarURL,
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
