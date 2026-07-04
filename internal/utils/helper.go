package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar,omitempty"`     // optional
	FullName  string    `json:"full_name,omitempty"`  // optional
	Gender    string    `json:"gender,omitempty"`     // optional
	CreatedAt time.Time `json:"created_at,omitempty"` // optional
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
		//Avatar:   user.AvatarURL,
		FullName:  user.Fullname,
		Gender:    user.Gender,
		CreatedAt: user.CreatedAt,
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func DeleteFileByPublicID(PublicID string) error {
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

func CalculateGrowth(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // or handle separately
	}
	return ((current - previous) / previous) * 100
}

func GenerateSlug(name string) string {
	return slug.Make(name)
}
