package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Quantity  int            `gorm:"not null;default:1"`
	AddedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
