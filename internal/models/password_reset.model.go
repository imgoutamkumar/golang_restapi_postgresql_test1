package models

import "time"

type PasswordReset struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       string    `gorm:"type:uuid;not null"`
	OTPHash      string    `gorm:"type:varchar(255);not null"`
	AttemptCount int       `gorm:"not null;default:0"`
	ExpiresAt    time.Time `gorm:"not null"`
	LockedAt     time.Time `gorm:""`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}
