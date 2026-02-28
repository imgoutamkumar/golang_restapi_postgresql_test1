package models

import "github.com/google/uuid"

type Wishlist struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string         `gorm:"type:uuid;not null;uniqueIndex"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name      string         `gorm:"type:varchar(100);not null;default:'My Wishlist'"`
	is_public bool           `gorm:"default:false"`
	Items     []WishlistItem `gorm:"foreignKey:WishlistID"`
}

type WishlistItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WishlistID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Wishlist       Wishlist       `gorm:"foreignKey:WishlistID;constraint:OnDelete:CASCADE"`
	VariantID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	ProductVariant ProductVariant `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`
}

