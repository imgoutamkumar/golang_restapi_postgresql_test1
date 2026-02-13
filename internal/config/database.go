package config

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DefaultUserRoleID uuid.UUID

func Connect(dsn string) (*gorm.DB, error) {
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// retry logic starts here
	var db *gorm.DB
	var err error

	for attempts := 1; attempts <= 10; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Println("dsn", dsn)
		log.Printf("DB connection attempt %d failed: %v", attempts, err)
		time.Sleep(2 * time.Second)
	}

	// retry logic ends here

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	DB = db
	log.Println("✅ Database connected")

	return DB, nil
}

func LoadDefaultRoles() {
	var role models.Role
	err := DB.Where("name = ?", "user").First(&role).Error

	if err != nil {
		log.Fatal("failed to load default user role:", err)
	}

	DefaultUserRoleID = role.ID
	if DefaultUserRoleID == uuid.Nil {
		log.Fatal("'user' role not found in roles table")
	}

	log.Println("Default user role loaded:", DefaultUserRoleID)

}
