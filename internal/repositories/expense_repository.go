package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
)

type ExpenseRepository struct{}

func NewExpenseRepository() *ExpenseRepository {
	return &ExpenseRepository{}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return db.DB.Create(expense).Error
}

func (r *ExpenseRepository) GetExpensesByUserID(userID string, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := db.DB.Where("user_id = ?", userID).Limit(limit).Order("ExpenseDate DESC").Find(&expenses).Error
	return expenses, err
}
