package main

import (
	"log"

	_ "visualfinance/docs"
	"visualfinance/internal/config"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/pkg/minio"
	"visualfinance/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Visual Finance API
// @version         1.0
// @description     This is the API for Visual Finance backend.
// @host            pem.danghoa-erp.site
// @BasePath        /api/v1
// @schemes         https

func main() {
	// 1. Load cấu hình từ file .env.dev
	config.LoadConfig()

	// 2. Khởi tạo kết nối CSDL (SQL Server)
	db.ConnectDB()

	// 2.5 Khởi tạo kết nối MinIO
	minio.ConnectMinIO()

	// 3. Thiết lập Gin Router
	router := gin.Default()

	// 3.5 Cấu hình CORS
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(configCors))

	// Khởi tạo các API routes
	routes.SetupExpenseRoutes(router)

	// Route cho Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 4. Chạy Server ở port 8080
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
