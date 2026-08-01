package models

import (
	"time"
)

type Session struct {
	ID           MSSQLUUID `gorm:"column:ID;type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID()" json:"id"`
	UserID       MSSQLUUID `gorm:"column:UserID;type:uniqueidentifier;not null" json:"user_id"`
	RefreshToken string    `gorm:"column:RefreshToken;type:varchar(500);not null;unique" json:"refresh_token"`
	UserAgent    string    `gorm:"column:UserAgent;type:varchar(255)" json:"user_agent"`
	ClientIP     string    `gorm:"column:ClientIP;type:varchar(50)" json:"client_ip"`
	IsBlocked    bool      `gorm:"column:IsBlocked;default:0" json:"is_blocked"`
	ExpiresAt    time.Time `gorm:"column:ExpiresAt;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:CreatedAt;default:GETDATE()" json:"created_at"`
}

func (Session) TableName() string {
	return "Sessions"
}


