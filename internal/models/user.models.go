package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Fullname  string    `gorm:"size:50;not null"`
	Username  string    `gorm:"size:50;unique;not null"`
	Gender    string    `gorm:"size:10;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RoleId    uuid.UUID `gorm:"type:uuid;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Orders []Order `gorm:"foreignKey:UserID"`
}
