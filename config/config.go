package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		return os.Getenv(key)
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, maybe docker trying to get environment variables")
	}
	return os.Getenv(key)
}
