package dto

import "github.com/google/uuid"

type TopSellingProductResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	TotalSold   int64     `json:"total_sold"`
}
