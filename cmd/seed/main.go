package main

import (
	"fmt"
	"log"

	"visualfinance/internal/config"
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Khởi tạo DB connection
	config.LoadConfig()
	db.ConnectDB()

	email := "dangquanghoa206@gmail.com"
	username := "danghoa"
	password := "admin123*"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsActive:     true, // Kích hoạt luôn không cần verify email
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}

	fmt.Println("User created successfully!")
}
