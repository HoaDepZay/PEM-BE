package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPRequest struct {
	ID        uuid.UUID `gorm:"type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID();column:ID"`
	Email     string    `gorm:"type:nvarchar(100);not null;column:Email"`
	OTPCode   string    `gorm:"type:varchar(6);not null;column:OTPCode"`
	ExpiresAt time.Time `gorm:"type:datetime2;not null;column:ExpiresAt"`
	IsUsed    bool      `gorm:"type:bit;default:0;column:IsUsed"`
	CreatedAt time.Time `gorm:"type:datetime2;default:GETDATE();column:CreatedAt"`
}

func (OTPRequest) TableName() string {
	return "OTPRequests"
}

// Hooks
func (o *OTPRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
