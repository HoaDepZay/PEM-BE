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

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	user, accessToken, refreshToken, err := authService.LoginUser(req.Email, req.Password, userAgent, clientIP)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Refresh issues a new access and refresh token
func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	newAccessToken, newRefreshToken, err := authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout deletes the refresh token session
func Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := authService.LogoutUser(req.RefreshToken); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
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
		"id":         user.UserID,
		"username":   user.Username,
		"email":      user.Email,
		"joinedAt":   user.CreatedAt,
		"avatar_url": user.AvatarURL,
	})
}

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Always return success to prevent email enumeration
	_ = authService.ForgotPassword(req.Email)
	response.Success(c, http.StatusOK, "Nếu email tồn tại trong hệ thống, mã OTP đã được gửi.", nil)
}

func VerifyOTP(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		OTPCode string `json:"otp_code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	resetToken, err := authService.VerifyOTP(req.Email, req.OTPCode)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OTP hợp lệ", gin.H{"reset_token": resetToken})
}

func ResetPassword(c *gin.Context) {
	var req struct {
		ResetToken  string `json:"reset_token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := authService.ResetPassword(req.ResetToken, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Đổi mật khẩu thành công. Bạn có thể đăng nhập.", nil)
}

func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(uuid.UUID)

	if err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Đổi mật khẩu thành công", nil)
}
