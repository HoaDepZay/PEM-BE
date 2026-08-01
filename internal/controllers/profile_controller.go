package controllers

import (
	"fmt"
	"net/http"

	"visualfinance/internal/pkg/minio"
	"visualfinance/internal/pkg/response"
	"visualfinance/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var profileService = services.NewProfileService()

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Avatar file is required")
		return
	}

	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(uuid.UUID)

	src, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer src.Close()

	// Generate unique filename
	fileName := fmt.Sprintf("%s-%s", userID.String(), file.Filename)

	// Upload to MinIO
	imageURL, err := minio.UploadImage(c.Request.Context(), fileName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload image")
		return
	}

	// Update User DB
	if err := profileService.UpdateAvatar(userID, imageURL); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update avatar in DB")
		return
	}

	response.Success(c, http.StatusOK, "Avatar uploaded successfully", gin.H{"avatar_url": imageURL})
}

func UpdateProfile(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input: username is required and must be at least 3 characters")
		return
	}

	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(uuid.UUID)

	if err := profileService.UpdateProfileInfo(userID, req.Username); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", nil)
}
