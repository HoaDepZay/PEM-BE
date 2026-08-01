package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	UserID       MSSQLUUID `gorm:"type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID();column:UserID"`
	Username     string    `gorm:"type:nvarchar(50);not null;unique;column:Username"`
	Email        string    `gorm:"type:nvarchar(100);not null;unique;column:Email"`
	PasswordHash string    `gorm:"type:nvarchar(255);not null;column:PasswordHash"`
	AvatarURL    string    `gorm:"type:varchar(500);column:AvatarURL" json:"avatar_url"`
	CreatedAt    time.Time `gorm:"type:datetime2;default:GETDATE();column:CreatedAt"`
	UpdatedAt    time.Time `gorm:"type:datetime2;default:GETDATE();column:UpdatedAt"`
	IsActive     bool      `gorm:"type:bit;default:0;column:IsActive"` // 0 for unverified, 1 for verified
}

func (User) TableName() string {
	return "Users"
}

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == NilMSSQLUUID {
		u.UserID = NewMSSQLUUID()
	}
	return
}



