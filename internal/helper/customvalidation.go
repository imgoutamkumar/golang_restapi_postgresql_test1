package helper

import (
	"errors"
)

// CustomValidate validates business rules for Product DTO
func CustomValidate(req *CreateProductRequest) error {
	// Default discount percent to 0
	discountPercent := req.DiscountPercent
	if discountPercent < 0 {
		discountPercent = 0
	}

	// Calculate discount amount
	discountAmount := (req.BasePrice * discountPercent) / 100

	// Discount must be less than base price
	if discountAmount >= req.BasePrice {
		return errors.New("discount price must be less than base price")
	}

	return nil
}
