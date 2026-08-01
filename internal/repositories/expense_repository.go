package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
	"gorm.io/gorm"
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
	err := db.DB.Where("UserID = ?", userID).Limit(limit).Order("ExpenseDate DESC").Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) GetByID(id string) (*models.Expense, error) {
	var expense models.Expense
	err := db.DB.First(&expense, "ExpenseID = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &expense, err
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return db.DB.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id string) error {
	return db.DB.Delete(&models.Expense{}, "ExpenseID = ?", id).Error
}

