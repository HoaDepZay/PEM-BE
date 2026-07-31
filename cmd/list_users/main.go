package main

import (
	"fmt"
	"visualfinance/internal/config"
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
)

func main() {
	config.LoadConfig()
	db.ConnectDB()

	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d users:\n", len(users))
	for _, u := range users {
		fmt.Printf("- Username: %s, Email: %s, IsActive: %t\n", u.Username, u.Email, u.IsActive)
	}
}
