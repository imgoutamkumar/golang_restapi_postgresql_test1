package services

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
)

func IsItemExists(wishlistId uuid.UUID, variantId uuid.UUID) (bool, error) {
	wishlistItem, err := repository.GetWishlistItemByVariantId(wishlistId, variantId)
	if err != nil && wishlistItem == nil {
		return false, err
	}

	return true, nil
}

func GetWishlistedVariantIds(userId uuid.UUID) ([]uuid.UUID, error) {
	wishlist, err := repository.GetWishlistByUserId(userId)
	if err != nil {
		return nil, err
	}
	variantIds := make([]uuid.UUID, len(wishlist.Items))
	for i, item := range wishlist.Items {
		variantIds[i] = item.VariantID
	}
	return variantIds, nil
}

func GetProductIdsFromvariantIds(variantIds []uuid.UUID) ([]uuid.UUID, error) {
	productIds, err := repository.GetProductIdsFromVariantIds(variantIds)
	if err != nil {
		return nil, err
	}
	return productIds, nil
}
