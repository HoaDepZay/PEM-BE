package main

import (
	"fmt"
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/repositories"
	"visualfinance/internal/pkg/utils"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	dsn := "sqlserver://sa:31052006Hoa*@100.109.65.2:1433?database=VisualFinanceDB"
	dialector := sqlserver.Open(dsn)
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}
	db.DB = gormDB

	repo := repositories.NewUserRepository()
	email := "dangquanghoa206@gmail.com"
	user, _ := repo.FindByEmail(email)
	
	fmt.Println("Original Active Status:", user.IsActive)
	user.IsActive = !user.IsActive
	
	err = db.DB.Model(&models.User{}).Where("UserID = ?", utils.ToMSSQLUUIDString(user.UserID)).Omit("UserID").Updates(user).Error
	if err != nil {
		fmt.Println("Update error:", err)
	} else {
		fmt.Println("Update Success without error")
	}
	
	user2, _ := repo.FindByEmail(email)
	fmt.Println("New Active Status:", user2.IsActive)
}
