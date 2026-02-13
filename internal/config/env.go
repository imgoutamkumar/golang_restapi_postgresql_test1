package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabaseUrl            string
	Port                   string
	Migrations             string
	JWTSecret              string
	CLOUDINARTY_API_KEY    string
	CLOUDINARTY_API_SECRET string
	CLOUDINARTY_CLOUD_NAME string
}

// public
func LoadEnv() *EnvConfig {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	var envConfig = &EnvConfig{
		DatabaseUrl:            os.Getenv("DB_URL"),
		Port:                   os.Getenv("PORT"),
		Migrations:             os.Getenv("MIGRATIONS_DIR"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
		CLOUDINARTY_API_KEY:    os.Getenv("CLOUDINARY_API_KEY"),
		CLOUDINARTY_API_SECRET: os.Getenv("CLOUDINARY_API_SECRET"),
		CLOUDINARTY_CLOUD_NAME: os.Getenv("CLOUDINARY_CLOUD_NAME"),
	}
	return envConfig
}

// private helper
func ValidateEnv(cfg *EnvConfig) error {
	// validation logic
	return nil
}
