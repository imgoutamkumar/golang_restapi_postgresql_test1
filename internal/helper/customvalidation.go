package helper

import (
	"errors"

	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
)

// CustomValidate validates business rules for Product DTO
func PriceValidate(req *dto.CreateProductRequest, i int) error {
	// Default discount percent to 0
	discountPercent := req.Variants[i].DiscountPercent
	if discountPercent < 0 {
		discountPercent = 0
	}

	// Calculate discount amount
	discountAmount := (req.Variants[i].Price * discountPercent) / 100

	// Discount must be less than base price
	if discountAmount >= req.Variants[i].Price {
		return errors.New("discount price must be less than base price")
	}

	return nil
}
