package main

import (
	"fmt"
	"log"
	"visualfinance/internal/config"
	"visualfinance/internal/models"
)

func main() {
	config.ConnectDB()
	var count int64
	if err := config.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total users in DB: %d\n", count)
}
