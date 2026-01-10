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
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	TotalAmount float64        `gorm:"type:numeric(10,2);not null"`
	Status      OrderStatus    `gorm:"type:order_status;default:'pending';not null"`
	CreatedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
