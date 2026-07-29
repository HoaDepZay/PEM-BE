package services

import (
	"context"
	"io"
	"time"

	"visualfinance/internal/models"
	"visualfinance/internal/pkg/minio"
	"visualfinance/internal/repositories"

	"github.com/google/uuid"
)

type ExpenseService struct {
	expenseRepo *repositories.ExpenseRepository
}

func NewExpenseService() *ExpenseService {
	return &ExpenseService{
		expenseRepo: repositories.NewExpenseRepository(),
	}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, userID, categoryID, note, filename, contentType string, amount, size int64, fileReader io.Reader) (*models.Expense, error) {
	// 1. Upload image to MinIO
	objectName := uuid.New().String() + "-" + filename
	imageURL, err := minio.UploadImage(ctx, objectName, fileReader, size, contentType)
	if err != nil {
		return nil, err
	}

	// 2. Create expense model
	expense := &models.Expense{
		UserID:      userID,
		Amount:      amount,
		Note:        note,
		ImageURL:    imageURL,
		ExpenseDate: time.Now(),
	}

	if categoryID != "" {
		expense.CategoryID = &categoryID
	}

	// 3. Save to DB
	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *ExpenseService) GetExpenses(userID string) ([]models.Expense, error) {
	return s.expenseRepo.GetExpensesByUserID(userID, 50)
}
