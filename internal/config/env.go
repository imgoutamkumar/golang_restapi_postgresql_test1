package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabaseUrl string
	Port        string
}

// public
func LoadEnv() (*EnvConfig, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("error")
	}
	var envConfig = &EnvConfig{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
	}
	return envConfig, err
}

// private helper
func validateEnv(cfg *EnvConfig) error {
	// validation logic
	return nil
}
