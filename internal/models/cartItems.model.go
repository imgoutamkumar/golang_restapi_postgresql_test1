package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItems struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Quantity  int            `gorm:"not null;default:1"`
	AddedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Cart    Cart    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
