package helper

import "github.com/shopspring/decimal"

type CreateProductRequest struct {
	Name            string          `json:"name" binding:"required,min=3,max=100"`
	Description     string          `json:"description" binding:"omitempty,max=2000"`
	BasePrice       decimal.Decimal `json:"price" binding:"required,gt=0"`
	DiscountPercent decimal.Decimal `json:"discount_percent"` // default 0 if not provided
}

type UserFilterParams struct {
	ProductName string
	FullName    string
	Page        int
	Limit       int
}
