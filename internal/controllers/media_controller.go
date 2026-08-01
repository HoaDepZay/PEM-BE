package controllers

import (
	"context"
	"net/http"

	"visualfinance/internal/pkg/minio"

	"github.com/gin-gonic/gin"
	minioV7 "github.com/minio/minio-go/v7"
)

func ServeMedia(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectName := c.Param("filename")
	
	// Khi dùng *filename, Gin có thể thêm / ở đầu
	if len(objectName) > 0 && objectName[0] == '/' {
		objectName = objectName[1:]
	}

	ctx := context.Background()
	object, err := minio.MinioClient.GetObject(ctx, bucketName, objectName, minioV7.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, nil)
}
