package services

import (
	"context"
	"errors"
	"io"
	"time"
	"fmt"

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

	uid, _ := uuid.Parse(userID)
	eid := uuid.New()

	// 2. Create expense model
	expense := &models.Expense{
		ExpenseID: models.MSSQLUUID(eid),
		UserID: models.MSSQLUUID(uid),
		Amount:      amount,
		Note:        note,
		ImageURL:    imageURL,
		ExpenseDate: time.Now(),
	}

	if categoryID != "" {
		cid, err := uuid.Parse(categoryID)
		if err == nil {
			expense.CategoryID = models.MSSQLUUID(cid)
		}
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

func (s *ExpenseService) UpdateExpense(userID, expenseID, categoryID, note, amountStr string) (*models.Expense, error) {
	expense, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, errors.New("giao dá»‹ch khÃ´ng tá»“n táº¡i")
	}
	uid, _ := uuid.Parse(userID)
	if uuid.UUID(expense.UserID) != uid {
		return nil, errors.New("báº¡n khÃ´ng cÃ³ quyá»n sá»­a giao dá»‹ch nÃ y")
	}

	if note != "" {
		expense.Note = note
	}
	if amountStr != "" {
		var amt int64
		fmt.Sscanf(amountStr, "%d", &amt)
		if amt > 0 {
			expense.Amount = amt
		}
	}
	if categoryID != "" {
		cid, err := uuid.Parse(categoryID)
		if err == nil {
			expense.CategoryID = models.MSSQLUUID(cid)
		}
	}

	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *ExpenseService) DeleteExpense(userID, expenseID string) error {
	expense, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return err
	}
	if expense == nil {
		return errors.New("giao dá»‹ch khÃ´ng tá»“n táº¡i")
	}
	uid, _ := uuid.Parse(userID)
	if uuid.UUID(expense.UserID) != uid {
		return errors.New("báº¡n khÃ´ng cÃ³ quyá»n xÃ³a giao dá»‹ch nÃ y")
	}

	return s.expenseRepo.Delete(expenseID)
}



