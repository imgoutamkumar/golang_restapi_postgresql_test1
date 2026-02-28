package dto

import "github.com/google/uuid"

type AddToWishlistRequest struct {
	VariantID uuid.UUID `json:"variant_id" validate:"required"`
}

type RemoveFromWishlistRequest struct {
	VariantID uuid.UUID `json:"variant_id" validate:"required"`
}
