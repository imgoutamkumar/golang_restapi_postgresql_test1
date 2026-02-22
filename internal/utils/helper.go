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
	Role     string    `json:"role"`
	Avatar   string    `json:"avatar,omitempty"` // optional
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
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
	varients := []dto.ProductVariantResponse{}
	for _, v := range p.Variants {
		variantImages := []dto.ProductImageResponse{}
		for _, img := range v.Images {
			variantImages = append(variantImages, dto.ProductImageResponse{
				URL:       img.ImageURL,
				IsPrimary: img.IsPrimary,
				PublicId:  img.PublicID,
			})
		}
		varients = append(varients, dto.ProductVariantResponse{
			Sku:             v.Sku,
			Price:           v.Price,
			DiscountPercent: v.DiscountPercent,
			FinalPrice:      v.Price - (v.Price * v.DiscountPercent / 100),
			Stock:           v.Stock,
			Images:          variantImages,
		})
	}

	return dto.ProductResponse{
		ID:        p.ID.String(),
		Name:      p.Name,
		ShortDesc: p.ShortDescription,
		Currency:  p.Currency,
		CreatedBy: p.CreatedBy.String(),
		CreatedAt: p.CreatedAt,
		Brand: dto.BrandResponse{
			ID:   p.BrandID.String(),
			Name: p.Brand.Name,
		},
	}
}

func GenerateCombinations(input [][]uuid.UUID) [][]uuid.UUID {
	if len(input) == 0 {
		return [][]uuid.UUID{}
	}

	result := [][]uuid.UUID{{}}

	for _, values := range input {
		var temp [][]uuid.UUID

		for _, r := range result {
			for _, v := range values {
				combination := append([]uuid.UUID{}, r...)
				combination = append(combination, v)
				temp = append(temp, combination)
			}
		}

		result = temp
	}

	return result
}
