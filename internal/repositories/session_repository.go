package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
)

type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return db.DB.Create(session).Error
}

func (r *SessionRepository) GetByRefreshToken(refreshToken string) (*models.Session, error) {
	var session models.Session
	err := db.DB.Where("RefreshToken = ?", refreshToken).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Delete(id string) error {
	return db.DB.Delete(&models.Session{}, "ID = ?", id).Error
}

func (r *SessionRepository) DeleteByToken(refreshToken string) error {
	return db.DB.Where("RefreshToken = ?", refreshToken).Delete(&models.Session{}).Error
}
