package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;index"`
	OrderNumber    string          `gorm:"type:varchar(30);not null;uniqueIndex"`
	Status         OrderStatus     `gorm:"type:order_status;not null;default:'pending'"`
	Subtotal       decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	DiscountAmount decimal.Decimal `gorm:"type:numeric(10,2);default:0"`
	TaxAmount      decimal.Decimal `gorm:"type:numeric(10,2);default:0"`
	ShippingAmount decimal.Decimal `gorm:"type:numeric(10,2);default:0"`
	TotalAmount    decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	CreatedAt      time.Time       `gorm:"not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	User       User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
