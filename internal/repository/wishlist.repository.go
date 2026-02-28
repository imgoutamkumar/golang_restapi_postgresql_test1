package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func CreateWishlist(wishlist *models.Wishlist) (*models.Wishlist, error) {
	err := config.DB.Create(wishlist).Error
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func GetDefaultWishlistByUserId(userId uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := config.DB.Where("user_id = ? AND name = ?", userId, "My Wishlist").First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func GetWishlistByUserId(userId uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := config.DB.Preload("Items").Where("user_id = ?", userId).First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func AddNewWishlistItem(wishlistItem *models.WishlistItem) error {
	return config.DB.Create(wishlistItem).Error
}

func GetWishlistItemByVariantId(wishlistId uuid.UUID, variantId uuid.UUID) (*models.WishlistItem, error) {
	var wishlistItem models.WishlistItem
	err := config.DB.Where("wishlist_id = ? AND variant_id = ?", wishlistId, variantId).First(&wishlistItem).Error
	if err != nil {
		return nil, err
	}
	return &wishlistItem, nil
}

func RemoveWishlistItem(wishlistId uuid.UUID, variantId uuid.UUID) error {
	return config.DB.Where("wishlist_id = ? AND variant_id = ?", wishlistId, variantId).Delete(&models.WishlistItem{}).Error
}

func RemoveProductFromAllWishlists(userId uuid.UUID, productId uuid.UUID) error {
	return config.DB.
		Table("wishlist_items wi").
		Joins("JOIN wishlists w ON w.id = wi.wishlist_id").
		Joins("JOIN product_variants pv ON pv.id = wi.variant_id").
		Where("w.user_id = ?", userId).
		Where("pv.product_id = ?", productId).
		Delete(nil).Error
}
