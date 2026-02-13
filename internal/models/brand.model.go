package models

import "github.com/google/uuid"

type Brand struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"size:50;not null"`
}
