package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
)

type OTPRepository struct{}

func NewOTPRepository() *OTPRepository {
	return &OTPRepository{}
}

func (r *OTPRepository) Create(otpReq *models.OTPRequest) error {
	return db.DB.Create(otpReq).Error
}

func (r *OTPRepository) GetValidOTP(email string, otpCode string) (*models.OTPRequest, error) {
	var otpReq models.OTPRequest
	// Need to check if it's unused and not expired
	err := db.DB.Where("Email = ? AND OTPCode = ? AND IsUsed = 0 AND ExpiresAt > GETDATE()", email, otpCode).First(&otpReq).Error
	if err != nil {
		return nil, err
	}
	return &otpReq, nil
}

func (r *OTPRepository) MarkAsUsed(id string) error {
	return db.DB.Model(&models.OTPRequest{}).Where("ID = ?", id).Update("IsUsed", 1).Error
}
