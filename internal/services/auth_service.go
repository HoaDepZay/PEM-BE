package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"visualfinance/internal/models"
	"visualfinance/internal/pkg/email"
	"visualfinance/internal/pkg/jwt"
	"visualfinance/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	otpRepo     *repositories.OTPRepository
	sessionRepo *repositories.SessionRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:    repositories.NewUserRepository(),
		otpRepo:     repositories.NewOTPRepository(),
		sessionRepo: repositories.NewSessionRepository(),
	}
}

func (s *AuthService) RegisterUser(username, emailStr, password string) (*models.User, string, error) {
	_, err := s.userRepo.FindByEmailOrUsername(emailStr, username)
	if err == nil {
		return nil, "", errors.New("Email or Username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("Failed to hash password")
	}

	newUser := &models.User{
		Username:     username,
		Email:        emailStr,
		PasswordHash: string(hashedPassword),
		IsActive:     false,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, "", errors.New("Failed to create user")
	}

	verifyToken, err := jwt.GenerateToken(uuid.UUID(newUser.UserID), newUser.Email, 24*time.Hour)
	if err != nil {
		return nil, "", errors.New("Failed to generate verification token")
	}

	return newUser, verifyToken, nil
}

func (s *AuthService) VerifyEmail(tokenString string) error {
	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return errors.New("Invalid or expired token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return errors.New("User not found")
	}

	if user.IsActive {
		return errors.New("Email is already verified")
	}

	user.IsActive = true
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("Failed to activate user")
	}

	return nil
}

func (s *AuthService) LoginUser(emailStr, password, userAgent, clientIP string) (*models.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(emailStr)
	if err != nil {
		return nil, "", "", errors.New("DB Error or User Not Found: " + err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", errors.New("Password mismatch")
	}

	accessToken, err := jwt.GenerateToken(uuid.UUID(user.UserID), user.Email, 15*time.Minute)
	if err != nil {
		return nil, "", "", errors.New("Failed to generate access token")
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", errors.New("Failed to generate refresh token")
	}

	session := &models.Session{
		ID: models.NewMSSQLUUID(),
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", "", errors.New("Failed to create session")
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.GetByRefreshToken(refreshToken)
	if err != nil || session == nil {
		return "", "", errors.New("Invalid refresh token")
	}

	if session.IsBlocked || time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.DeleteByToken(refreshToken)
		return "", "", errors.New("Refresh token expired or blocked")
	}

	user, err := s.userRepo.FindByID(uuid.UUID(session.UserID))
	if err != nil {
		return "", "", errors.New("User not found")
	}

	newAccessToken, err := jwt.GenerateToken(uuid.UUID(user.UserID), user.Email, 15*time.Minute)
	if err != nil {
		return "", "", errors.New("Failed to generate new access token")
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	// XÃ³a session cÅ©, táº¡o session má»›i (Xoay vÃ²ng token)
	_ = s.sessionRepo.DeleteByToken(refreshToken)
	
	newSession := &models.Session{
		ID: models.NewMSSQLUUID(),
		UserID:       user.UserID,
		RefreshToken: newRefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	_ = s.sessionRepo.Create(newSession)

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) LogoutUser(refreshToken string) error {
	return s.sessionRepo.DeleteByToken(refreshToken)
}


func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (s *AuthService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// Generate random 6-digit OTP
func generateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func (s *AuthService) ForgotPassword(emailStr string) error {
	user, err := s.userRepo.FindByEmail(emailStr)
	if err != nil || user == nil {
		// Do not leak if email exists or not for security, just return nil
		return nil
	}

	otpCode := generateOTP()
	otpReq := &models.OTPRequest{
		Email:     emailStr,
		OTPCode:   otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.otpRepo.Create(otpReq); err != nil {
		return errors.New("Failed to process request")
	}

	// Send OTP email in background
	go func() {
		err := email.SendOTPEmail(emailStr, otpCode)
		if err != nil {
			log.Printf("Error sending OTP email: %v", err)
		} else {
			log.Printf("Successfully sent OTP email to %s", emailStr)
		}
	}()

	return nil
}

func (s *AuthService) VerifyOTP(emailStr, otpCode string) (string, error) {
	otpReq, err := s.otpRepo.GetValidOTP(emailStr, otpCode)
	if err != nil || otpReq == nil {
		return "", errors.New("Invalid or expired OTP")
	}

	// Mark as used
	_ = s.otpRepo.MarkAsUsed(otpReq.ID.String())

	// Generate reset token
	user, err := s.userRepo.FindByEmail(emailStr)
	if err != nil {
		return "", errors.New("User not found")
	}

	resetToken, err := jwt.GenerateToken(uuid.UUID(user.UserID), user.Email, 15*time.Minute)
	if err != nil {
		return "", errors.New("Failed to generate reset token")
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(resetToken, newPassword string) error {
	claims, err := jwt.ValidateToken(resetToken)
	if err != nil {
		return errors.New("Invalid or expired reset token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return errors.New("User not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("Failed to update password")
	}

	return nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("Incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash new password")
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("Failed to change password")
	}

	return nil
}


