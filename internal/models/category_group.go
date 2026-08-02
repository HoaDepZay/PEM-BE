package models

import (
	"time"
)

type CategoryGroup struct {
	GroupID     MSSQLUUID `gorm:"column:GroupID;type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID()" json:"group_id"`
	UserID      *MSSQLUUID `gorm:"column:UserID;type:uniqueidentifier" json:"user_id"` // NULL means global/system group
	Name        string    `gorm:"column:Name;type:nvarchar(100)" json:"name"`
	Icon        string    `gorm:"column:Icon;type:varchar(50)" json:"icon"`
	Color       string    `gorm:"column:Color;type:varchar(20)" json:"color"`
	TotalBudget int64     `gorm:"column:TotalBudget;type:bigint;default:0" json:"total_budget"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;default:GETDATE()" json:"created_at"`

	// Relationships
	Categories []Category `gorm:"foreignKey:GroupID;references:GroupID" json:"categories,omitempty"`
}

func (CategoryGroup) TableName() string {
	return "CategoryGroups"
}
