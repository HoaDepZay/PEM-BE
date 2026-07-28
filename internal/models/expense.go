package models

import (
	"time"
)

type Expense struct {
	ExpenseID   string    `gorm:"column:ExpenseID;type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID()" json:"expense_id"`
	UserID      string    `gorm:"column:UserID;type:uniqueidentifier;not null" json:"user_id"`
	CategoryID  *string   `gorm:"column:CategoryID;type:uniqueidentifier" json:"category_id"`
	Amount      int64     `gorm:"column:Amount;not null" json:"amount"`
	Note        string    `gorm:"column:Note;type:nvarchar(255)" json:"note"`
	ImageURL    string    `gorm:"column:ImageURL;type:varchar(500)" json:"image_url"`
	ExpenseDate time.Time `gorm:"column:ExpenseDate;not null" json:"expense_date"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;default:GETDATE()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt;default:GETDATE()" json:"updated_at"`
	IsDeleted   bool      `gorm:"column:IsDeleted;default:0" json:"is_deleted"`
}

// TableName overrides the table name used by Expense to `Expenses`
func (Expense) TableName() string {
	return "Expenses"
}
