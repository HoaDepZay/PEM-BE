package main

import (
	"log"

	"visualfinance/internal/config"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load cấu hình từ file .env.dev
	config.LoadConfig()

	// 2. Khởi tạo kết nối CSDL (SQL Server)
	db.ConnectDB()

	// 3. Thiết lập Gin Router
	router := gin.Default()

	// Khởi tạo các API routes
	routes.SetupExpenseRoutes(router)

	// 4. Chạy Server ở port 8080
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
