package services

import (
	"errors"

	"visualfinance/internal/models"
	"visualfinance/internal/repositories"

	"github.com/google/uuid"
)

type ProfileService struct {
	userRepo *repositories.UserRepository
}

func NewProfileService() *ProfileService {
	return &ProfileService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *ProfileService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (s *ProfileService) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	user.AvatarURL = avatarURL
	return s.userRepo.Update(user)
}

func (s *ProfileService) UpdateProfileInfo(userID uuid.UUID, username string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	// Optional check: Ensure the new username isn't already taken by someone else
	if username != user.Username {
		existingUser, err := s.userRepo.FindByEmailOrUsername("", username)
		if err == nil && existingUser.UserID != user.UserID {
			return errors.New("Username is already taken")
		}
	}

	user.Username = username
	return s.userRepo.Update(user)
}
