package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	// Giả sử chạy từ backend/cmd/api
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Println("Could not load .env.dev file, falling back to system environment variables")
	}
}
