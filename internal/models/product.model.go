package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:numeric(10,2);not null"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	User      User      `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CartItems []Cart `gorm:"foreignKey:ProductID"`
}
