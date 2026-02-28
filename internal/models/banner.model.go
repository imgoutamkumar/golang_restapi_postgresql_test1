package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Banner struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:varchar(50);not null"` // hero, sale, promo
	IsActive  bool      `gorm:"default:true"`
	StartDate *time.Time
	EndDate   *time.Time

	Images []BannerImage `gorm:"foreignKey:BannerID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type BannerImage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BannerID uuid.UUID `gorm:"type:uuid;not null"`

	ImageURL string `gorm:"type:text;not null"`
	PublicID string `gorm:"type:varchar(255);not null"`

	LinkURL     string `gorm:"type:text"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`

	IsActive  bool `gorm:"default:true"`
	SortOrder int  `gorm:"default:0"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Banner Banner `gorm:"foreignKey:BannerID"`
}
