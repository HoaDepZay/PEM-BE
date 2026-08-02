package models

import (
	"time"
)

type Category struct {
	CategoryID   MSSQLUUID `gorm:"column:CategoryID;type:uniqueidentifier;primaryKey;default:NEWSEQUENTIALID()" json:"category_id"`
	GroupID      MSSQLUUID `gorm:"column:GroupID;type:uniqueidentifier;not null" json:"group_id"`
	Name         string    `gorm:"column:Name;type:nvarchar(100)" json:"name"`
	BudgetType   string    `gorm:"column:BudgetType;type:varchar(20)" json:"budget_type"`
	BudgetAmount int64     `gorm:"column:BudgetAmount;type:bigint" json:"budget_amount"`
	CreatedAt    time.Time `gorm:"column:CreatedAt;default:GETDATE()" json:"created_at"`

	// Virtual fields for dynamic calculation (Not in DB)
	Status           string `gorm:"-" json:"status"`
	OverBudgetAmount int64  `gorm:"-" json:"over_budget_amount"`
}

func (Category) TableName() string {
	return "Categories"
}
