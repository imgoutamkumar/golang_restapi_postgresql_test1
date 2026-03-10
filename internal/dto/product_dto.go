package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name             string `form:"name" validate:"required,min=3"`
	BrandID          string `form:"brand_id" validate:"required"`
	CategoryID       string `form:"category_id" validate:"required"`
	ShortDescription string `form:"short_description"`
	Description      string `form:"description"`
	Currency         string `form:"currency"`
	Status           string `form:"status" validate:"required,oneof=active inactive draft"`
	IsReturnable     bool   `form:"is_returnable"`
	IsCodAvailable   bool   `form:"is_cod_available"`
	// Variants    []CreateVariantRequest `form:"variants" validate:"required,dive"`
	VariantsJSON string `form:"variants"`
}

type CreateVariantRequest struct {
	Sku               string   `form:"sku" validate:"required"`
	Price             float64  `form:"price" validate:"required,gt=0"`
	Stock             int      `form:"stock" validate:"gte=0"`
	DiscountPercent   float64  `form:"discount_percent" validate:"gte=0,lte=100"`
	IsDefault         bool     `form:"is_default"`
	Status            string   `form:"status" validate:"required,oneof=active inactive"`
	AttributeValueIDs []string `form:"attribute_value_ids" validate:"required"`

	// images for this variant
	ImageFiles   []*multipart.FileHeader `form:"images" validate:"required"`
	PrimaryIndex int                     `form:"primary_index"` // which image is primary
}

type CreateImageRequest struct {
	ImageURL  string
	PublicID  string
	IsPrimary bool
	SortOrder int
}

type ProductImageResponse struct {
	Id        string `json:"id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
	PublicId  string `json:"public_id"`
}

type BrandResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AttributeTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type AttributeValueResponse struct {
	ID       string `json:"id"`
	Value    string `json:"value"`
	MetaInfo string `json:"meta_info"`
}

type AttributeGroup struct {
	Name   string
	Values []AttributeValueResponse
}

type ProductVariantResponse struct {
	Id              string                 `json:"id"`
	Sku             string                 `json:"sku"`
	Price           float64                `json:"price"`
	DiscountPercent float64                `json:"discount_percent"`
	OriginalPrice   float64                `json:"original_price"`
	FinalPrice      float64                `json:"final_price"`
	Stock           int                    `json:"stock"`
	Status          string                 `jdon:"status"`
	AttributeGroup  []AttributeGroup       `json:"attribute_groups"`
	Images          []ProductImageResponse `json:"images"`
}

type ProductResponse struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	ShortDesc    string                   `json:"short_description"`
	BasePrice    float64                  `json:"base_price"`
	Currency     string                   `json:"currency"`
	Brand        BrandResponse            `json:"brand"`
	CreatedBy    string                   `json:"created_by"`
	Variants     []ProductVariantResponse `json:"variants"`
	CreatedAt    time.Time                `json:"created_at"`
	IsWishlisted bool                     `json:"is_wishlisted"`
}
type FiltersResponse struct {
	Brands []BrandResponse `json:"brands"`
}

type CreateBrandRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type CreateCategoryRequest struct {
	Name     string     `json:"name" validate:"required,min=2"`
	ParentID *uuid.UUID `json:"parent_id" validate:"omitempty,uuid"`
}

type CreateAttributeRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type CreateAttributeValueRequest struct {
	Value    string `json:"value" validate:"required,min=2"`
	MetaInfo string `json:"meta_info"`
	// AttributeTypeID uuid.UUID `json:"attribute_type_id" validate:"required"`
}

type ReorderProductImagesRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`

	Images []struct {
		ID        string `json:"id" binding:"required,uuid"`
		SortOrder int    `json:"sort_order" binding:"gte=0"`
	} `json:"images" binding:"required,min=1"`
}
