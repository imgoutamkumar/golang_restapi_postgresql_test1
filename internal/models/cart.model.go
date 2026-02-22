package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;unique;index"`
	AddedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CartItems []CartItems `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	User      User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type CartItems struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	VariantID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Quantity  int            `gorm:"not null;default:1"`
	AddedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Cart    Cart           `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Variant ProductVariant `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`
}
