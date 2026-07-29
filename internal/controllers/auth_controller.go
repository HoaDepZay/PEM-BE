package controllers

import (
	"fmt"
	"net/http"

	"visualfinance/internal/pkg/email"
	"visualfinance/internal/pkg/response"
	"visualfinance/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var authService = services.NewAuthService()

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user account and sends a verification email
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	newUser, verifyToken, err := authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	frontendURL := "http://localhost:5173"
	if gin.Mode() == gin.ReleaseMode {
		frontendURL = "https://pem.danghoa-erp.site"
	}
	verifyLink := fmt.Sprintf("%s/verify?token=%s", frontendURL, verifyToken)

	// Send Email
	go func() {
		if err := email.SendVerificationEmail(newUser.Email, verifyLink); err != nil {
			fmt.Println("Error sending email:", err)
		}
	}()

	response.Success(c, http.StatusOK, "Registration successful. Please check your email to verify your account.", nil)
}

// VerifyEmail verifies the user's email via the token
func VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		response.Error(c, http.StatusBadRequest, "Token is required")
		return
	}

	if err := authService.VerifyEmail(tokenString); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Email verified successfully. You can now login.", nil)
}

// Login authenticates a user and returns an access token
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	user, accessToken, err := authService.LoginUser(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": accessToken,
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetMe returns the current authenticated user's profile
func GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := userIDStr.(uuid.UUID)
	user, err := authService.GetUserByID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved", gin.H{
		"id":       user.UserID,
		"username": user.Username,
		"email":    user.Email,
		"joinedAt": user.CreatedAt,
	})
}
