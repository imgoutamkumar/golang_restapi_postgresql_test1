package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductDraft    ProductStatus = "draft"
	ProductActive   ProductStatus = "active"
	ProductInactive ProductStatus = "inactive"
	ProductArchived ProductStatus = "archived"
)

type Brand struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string
}

type Category struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string
}

type AttributeType struct {
	ID     uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name   string           `gorm:"type:varchar(50);uniqueIndex;not null"`
	Values []AttributeValue `gorm:"foreignKey:AttributeTypeID;constraint:OnDelete:CASCADE"`
}

type AttributeValue struct {
	ID              uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AttributeTypeID uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex:idx_type_value"`
	Value           string        `gorm:"type:varchar(100);not null;uniqueIndex:idx_type_value"`
	MetaInfo        string        `gorm:"type:varchar(100)"` // e.g., "#FF0000" for Red
	AttributeType   AttributeType `gorm:"foreignKey:AttributeTypeID"`
}

type Product struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name             string         `gorm:"type:varchar(150);not null;index"`
	BrandID          uuid.UUID      `gorm:"type:uuid;not null;index"`
	CategoryID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	Description      string         `gorm:"type:text"`
	ShortDescription string         `gorm:"type:varchar(500)"`
	BasePrice        float64        `gorm:"type:numeric(10,2);not null"`
	Currency         string         `gorm:"type:char(3);default:'INR'"`
	Status           ProductStatus  `gorm:"type:varchar(50);default:'draft';index"` // Assuming product_status maps to a string in Go
	IsReturnable     bool           `gorm:"default:true"`
	IsCodAvailable   bool           `gorm:"default:true"`
	CreatedBy        uuid.UUID      `gorm:"type:uuid;not null"`
	Details          datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	Slug             string         `gorm:"type:varchar(160);uniqueIndex"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	Brand        Brand            `gorm:"foreignKey:BrandID"`
	Category     Category         `gorm:"foreignKey:CategoryID"`
	Variants     []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	AvgRating    float64          `gorm:"type:numeric(3,2);default:0"`
	TotalReviews int              `gorm:"default:0"`
}

type ProductVariant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`

	Sku             string        `gorm:"type:varchar(100);uniqueIndex;not null"`
	Price           float64       `gorm:"type:numeric(10,2);not null"`
	DiscountPercent float64       `gorm:"type:numeric(5,2);default:0;check:discount_percent >= 0 AND discount_percent <= 100"`
	IsDefault       bool          `gorm:"default:false"`
	Stock           int           `gorm:"not null;default:0;check:stock >= 0"`
	Status          ProductStatus `gorm:"type:varchar(20);default:'active';index"`
	Slug            string        `gorm:"type:varchar(120);uniqueIndex"`
	CreatedAt       time.Time
	Product         Product `gorm:"constraint:OnDelete:CASCADE;"`

	Images []ProductImage `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`

	VariantAttributes []VariantAttribute `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VariantID uuid.UUID `gorm:"type:uuid;not null;index"`
	ImageURL  string    `gorm:"type:text;not null"`
	PublicID  string    `gorm:"type:varchar(255);not null"`
	IsPrimary bool      `gorm:"default:false"`
	SortOrder int       `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type VariantAttribute struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VariantID        uuid.UUID `gorm:"type:uuid;not null"`
	AttributeValueID uuid.UUID `gorm:"type:uuid;not null"`

	Variant        ProductVariant
	AttributeValue AttributeValue
}
