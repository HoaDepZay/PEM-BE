package main

import (
	"fmt"
	"time"

	"visualfinance/internal/pkg/jwt"
	"github.com/google/uuid"
)

func main() {
	id := uuid.New()
	email := "test@example.com"
	
	fmt.Println("Original ID:", id)

	token, err := jwt.GenerateToken(id, email, 15*time.Minute)
	if err != nil {
		fmt.Println("Error generating:", err)
		return
	}

	fmt.Println("Token:", token)

	claims, err := jwt.ValidateToken(token)
	if err != nil {
		fmt.Println("Error validating:", err)
		return
	}

	fmt.Println("Decoded ID:", claims.UserID)
	fmt.Println("Decoded Email:", claims.Email)
}
