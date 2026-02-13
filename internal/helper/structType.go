package helper

import (
	"mime/multipart"

	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

type ProductImageRequest struct {
	ImageURL  string `json:"imageUrl" validate:"required,url"`
	IsPrimary bool   `json:"isPrimary"`
	SortOrder int    `json:"sortOrder"`
}

// type CreateProductRequest struct {
// 	Name            string                  `json:"name" binding:"required,min=3,max=100"`
// 	Description     string                  `json:"description" binding:"omitempty,max=2000"`
// 	BasePrice       decimal.Decimal         `json:"price" binding:"required,gt=0"`
// 	DiscountPercent decimal.Decimal         `json:"discount_percent"`
// }

// if you are using form-data for file upload
type CreateProductRequest struct {
	Name             string               `form:"productname" validate:"required,min=3"`
	BrandID          string               `form:"brandId" validate:"required,uuid"`
	Category         string               `form:"category" validate:"required"`
	BasePrice        float64              `form:"base_price" validate:"gt=0"`
	NumberOfStock    int                  `form:"number_of_stock" validate:"gte=0"`
	DiscountPercent  float64              `form:"discount_percent"`
	Currency         string               `form:"currency"`
	Status           models.ProductStatus `form:"status" validate:"required,oneof=draft active inactive archived"`
	IsReturnable     bool                 `form:"is_returnable"`
	IsCODAvailable   bool                 `form:"is_cod_available"`
	Description      string               `form:"description"`
	ShortDescription string               `form:"short_description"`

	ImageFiles   []*multipart.FileHeader `form:"product_images"`
	PrimaryIndex int                     `form:"primary_index"`
}
type UserFilterParams struct {
	ProductName string
	FullName    string
	Page        int
	Limit       int
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email validate:"required,email"`
}

type VerifyOtpRequestBody struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required,len=6"`
}

type ReorderProductImagesRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`

	Images []struct {
		ID        string `json:"id" binding:"required,uuid"`
		SortOrder int    `json:"sort_order" binding:"gte=0"`
	} `json:"images" binding:"required,min=1"`
}
