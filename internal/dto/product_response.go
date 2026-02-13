package dto

import "time"

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

type ProductResponse struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	ShortDesc       string                 `json:"short_description"`
	BasePrice       float64                `json:"price"`
	DiscountPercent float64                `json:"discount_percent"`
	FinalPrice      float64                `json:"final_price"`
	Currency        string                 `json:"currency"`
	Stock           int                    `json:"stock"`
	Brand           BrandResponse          `json:"brand"`
	CreatedBy       string                 `json:"created_by"`
	Images          []ProductImageResponse `json:"images"`
	CreatedAt       time.Time              `json:"created_at"`
}
