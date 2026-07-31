package main

import (
	"fmt"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/repositories"
	"visualfinance/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func MSSQLUUIDString(u uuid.UUID) string {
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		u[3], u[2], u[1], u[0],
		u[5], u[4],
		u[7], u[6],
		u[8], u[9],
		u[10], u[11], u[12], u[13], u[14], u[15])
}

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
	
	fmt.Println("Go UUID string:", user.UserID.String())
	fmt.Println("Converted SQL string:", MSSQLUUIDString(user.UserID))

	var u2 models.User
	err = db.DB.Where("UserID = ?", MSSQLUUIDString(user.UserID)).First(&u2).Error
	if err != nil {
		fmt.Println("Query by converted string failed:", err)
	} else {
		fmt.Println("Query by converted string SUCCESS:", u2.Email)
	}
}
