package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderShipped   OrderStatus = "shipped"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index"`
	OrderNumber    string         `gorm:"type:varchar(30);not null;uniqueIndex"`
	Status         OrderStatus    `gorm:"type:order_status;not null;default:'pending'"`
	Subtotal       float64        `gorm:"type:numeric(10,2);not null"`
	DiscountAmount float64        `gorm:"type:numeric(10,2);default:0"`
	TaxAmount      float64        `gorm:"type:numeric(10,2);default:0"`
	ShippingAmount float64        `gorm:"type:numeric(10,2);default:0"`
	TotalAmount    float64        `gorm:"type:numeric(10,2);not null"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	User       User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	VariantID uuid.UUID `gorm:"type:uuid;not null;index"`

	ProductName string `gorm:"type:varchar(150);not null"`

	// snapshot at purchase time
	ProductPrice    float64 `gorm:"type:numeric(10,2);not null"`
	DiscountPercent float64 `gorm:"type:numeric(10,2);default:0"`
	Quantity        int     `gorm:"not null;check:quantity > 0"`
	TotalPrice      float64 `gorm:"type:numeric(10,2);not null"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// Relations
	Order   Order          `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Variant ProductVariant `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`

	// Unique constraint
	// uq_order_product (order_id, product_id)
}
