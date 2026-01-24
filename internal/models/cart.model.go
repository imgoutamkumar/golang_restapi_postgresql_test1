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
