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
	RoleID    uuid.UUID `gorm:"type:uuid;not null"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Orders []Order `gorm:"foreignKey:UserID"`
}

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"size:50;not null"`
}

type PasswordReset struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       string    `gorm:"type:uuid;not null"`
	OTPHash      string    `gorm:"type:varchar(255);not null"`
	AttemptCount int       `gorm:"not null;default:0"`
	ExpiresAt    time.Time `gorm:"not null"`
	LockedAt     time.Time `gorm:""`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}
