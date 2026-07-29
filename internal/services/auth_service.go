package services

import (
	"errors"
	"time"

	"visualfinance/internal/models"
	"visualfinance/internal/pkg/jwt"
	"visualfinance/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, string, error) {
	// 1. Check if exists
	_, err := s.userRepo.FindByEmailOrUsername(email, username)
	if err == nil {
		return nil, "", errors.New("Email or Username already exists")
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("Failed to hash password")
	}

	// 3. Create user
	newUser := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsActive:     false,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, "", errors.New("Failed to create user")
	}

	// 4. Generate verification token
	verifyToken, err := jwt.GenerateToken(newUser.UserID, newUser.Email, 24*time.Hour)
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

func (s *AuthService) LoginUser(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("Invalid email or password")
	}

	if !user.IsActive {
		return nil, "", errors.New("Please verify your email before logging in")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("Invalid email or password")
	}

	accessToken, err := jwt.GenerateToken(user.UserID, user.Email, 7*24*time.Hour)
	if err != nil {
		return nil, "", errors.New("Failed to generate token")
	}

	return user, accessToken, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
