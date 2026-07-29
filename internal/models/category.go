package models

import (
	"time"
)

type Category struct {
	CategoryID string    `gorm:"column:CategoryID;type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID()" json:"category_id"`
	UserID     *string   `gorm:"column:UserID;type:uniqueidentifier" json:"user_id"` // NULL means global/system category
	Name       string    `gorm:"column:Name;type:nvarchar(100)" json:"name"`
	Icon       string    `gorm:"column:Icon;type:varchar(50)" json:"icon"`
	Color      string    `gorm:"column:Color;type:varchar(20)" json:"color"`
	CreatedAt  time.Time `gorm:"column:CreatedAt;default:GETDATE()" json:"created_at"`
}

// TableName overrides the table name used by Category to `Categories`
func (Category) TableName() string {
	return "Categories"
}
