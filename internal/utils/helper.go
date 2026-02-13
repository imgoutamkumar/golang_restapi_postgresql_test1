package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
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

func DeleteFile(fileURL string) error {
	// OPTION A: Delete from Local Disk
	// path := "." + fileURL // assuming fileURL is like "/uploads/filename.jpg"
	// return os.Remove(path)

	// OPTION B: Delete from AWS S3
	// This is where you would use the AWS SDK to delete the file from the bucket.
	// For this example, I will simulate a successful deletion.
	// Simulate processing time
	// time.Sleep(50 * time.Millisecond)
	return nil
}

func GenerateOtp() string {
	// Generate a random 6-digit OTP
	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	return otp
}

func MapProductToResponse(p *models.Product) dto.ProductResponse {
	finalPrice := p.BasePrice - (p.BasePrice * p.DiscountPercent / 100)

	images := []dto.ProductImageResponse{}
	for _, img := range p.ProductImages {
		images = append(images, dto.ProductImageResponse{
			URL:       img.ImageUrl,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicId,
		})
	}

	return dto.ProductResponse{
		ID:              p.ID.String(),
		Name:            p.Name,
		ShortDesc:       p.ShortDescription,
		BasePrice:       p.BasePrice,
		DiscountPercent: p.DiscountPercent,
		FinalPrice:      finalPrice,
		Currency:        p.Currency,
		Stock:           p.NumberOfStock,
		Brand: dto.BrandResponse{
			ID:   p.Brand.ID.String(),
			Name: p.Brand.Name,
		},
		Images:    images,
		CreatedAt: p.CreatedAt,
	}
}
