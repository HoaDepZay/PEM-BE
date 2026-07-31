package main

import (
	"fmt"
	"visualfinance/internal/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	dsn := "sqlserver://sa:31052006Hoa*@100.109.65.2:1433?database=VisualFinanceDB"
	dialector := sqlserver.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}

	var users []models.User
	db.Find(&users)
	for _, u := range users {
		fmt.Printf("User: %s, ID: %v\n", u.Email, u.UserID)
	}

	var rawID string
	db.Raw("SELECT CAST(UserID AS VARCHAR(50)) FROM Users WHERE Email = 'dangquanghoa206@gmail.com'").Scan(&rawID)
	fmt.Println("Raw SQL ID:", rawID)
}
