package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductDraft    ProductStatus = "draft"
	ProductActive   ProductStatus = "active"
	ProductInactive ProductStatus = "inactive"
	ProductArchived ProductStatus = "archived"
)

type Product struct {
	ID               uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name             string        `gorm:"size:100;not null"`
	ShortDescription string        `gorm:"type:text"`
	BasePrice        float64       `gorm:"type:numeric(10,2);not null"`
	DiscountPercent  float64       `gorm:"type:numeric(10,2);default:0;check:discount_percent >= 0 AND discount_percent <= 100"`
	Currency         string        `gorm:"type:char(3);default:'INR'"`
	IsReturnable     bool          `gorm:"type:boolean;default:true"`
	IsCodAvailable   bool          `gorm:"type:boolean;default:true"`
	NumberOfStock    int           `gorm:"type:integer;default:0;check:number_of_stock >= 0"`
	Status           ProductStatus `gorm:"type:product_status;default:'draft'"`
	Description      string        `gorm:"type:text"`
	IsCODAvailable   bool          `gorm:"type:boolean;default:true"`
	BrandID          uuid.UUID     `gorm:"size:100"`
	Brand            Brand         `gorm:"foreignKey:BrandID"`
	// Category         string        `gorm:"type:uuid;not null"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	User      User      `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ProductImages []ProductImages `gorm:"foreignKey:ProductID"`
}
