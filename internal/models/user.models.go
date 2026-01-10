package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string         `gorm:"size:50;unique;not null"`
	Email     string         `gorm:"size:100;unique;not null"`
	Password  string         `gorm:"size:255;not null"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Orders []Order `gorm:"foreignKey:UserID"`
	Cart   []Cart  `gorm:"foreignKey:UserID"`
}
