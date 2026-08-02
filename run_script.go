package main

import (
	"log"
	"visualfinance/internal/config"
	"visualfinance/internal/pkg/db"
)

func main() {
	config.LoadConfig()
	db.ConnectDB()

	// Drop trigger since GORM doesn't support OUTPUT clause with triggers in SQL Server
	if err := db.DB.Exec("DROP TRIGGER IF EXISTS trg_UpdateTotalBudget").Error; err != nil {
		log.Printf("Failed to drop trigger: %v\n", err)
	} else {
		log.Println("Successfully dropped trigger")
	}
}
