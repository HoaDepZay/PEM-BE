package main

import (
	"fmt"
	"log"

	"visualfinance/internal/pkg/db"
	"visualfinance/internal/repositories"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()

	repo := repositories.NewUserRepository()
	user, err := repo.FindByEmail("dangquanghoa206@gmail.com")
	if err != nil {
		log.Fatalf("Error finding user: %v", err)
	}

	fmt.Printf("User ID: %s, Email: %s\n", user.UserID, user.Email)
	fmt.Printf("IsActive: %v\n", user.IsActive)
	fmt.Printf("PasswordHash: %s\n", user.PasswordHash)

	passwordToTest := "Admin123*"
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordToTest))
	if err != nil {
		fmt.Printf("Password %s does NOT match! Error: %v\n", passwordToTest, err)
	} else {
		fmt.Printf("Password %s matches!\n", passwordToTest)
	}
}
