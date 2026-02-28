package helper

// func PriceValidate(req *dto.CreateProductRequest, i int) error {
// 	discountPercent := req.Variants[i].DiscountPercent
// 	if discountPercent < 0 {
// 		discountPercent = 0
// 	}
// 	discountAmount := (req.Variants[i].Price * discountPercent) / 100
// 	if discountAmount >= req.Variants[i].Price {
// 		return errors.New("discount price must be less than base price")
// 	}
// 	return nil
// }
