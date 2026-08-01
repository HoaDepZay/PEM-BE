package services

import (
	"errors"
	"visualfinance/internal/models"
	"visualfinance/internal/repositories"
	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repositories.NewCategoryRepository(),
	}
}

func (s *CategoryService) GetCategories(userID string) ([]models.Category, error) {
	return s.categoryRepo.GetCategoriesByUserID(userID)
}

func (s *CategoryService) CreateCategory(userID, name, icon, color string, dailyBudget *int64) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("tÃªn danh má»¥c khÃ´ng Ä‘Æ°á»£c Ä‘á»ƒ trá»‘ng")
	}

	uid, _ := uuid.Parse(userID)

	category := &models.Category{
		CategoryID: models.NewMSSQLUUID(),
		UserID: (*models.MSSQLUUID)(&uid),
		Name:   name,
		Icon:   icon,
		Color:       color,
		DailyBudget: dailyBudget,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(userID, categoryID, name, icon, color string, dailyBudget *int64) (*models.Category, error) {
	// Kiá»ƒm tra xem category cÃ³ tá»“n táº¡i vÃ  thuá»™c vá» user khÃ´ng
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("danh má»¥c khÃ´ng tá»“n táº¡i")
	}
	uid, _ := uuid.Parse(userID)
	if category.UserID == nil || uuid.UUID(*category.UserID) != uid {
		return nil, errors.New("báº¡n khÃ´ng cÃ³ quyá»n sá»­a danh má»¥c nÃ y")
	}

	if name != "" {
		category.Name = name
	}
	if icon != "" {
		category.Icon = icon
	}
	if color != "" {
		category.Color = color
	}
	category.DailyBudget = dailyBudget

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(userID, categoryID string) error {
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("danh má»¥c khÃ´ng tá»“n táº¡i")
	}
	uid, _ := uuid.Parse(userID)
	if category.UserID == nil || uuid.UUID(*category.UserID) != uid {
		return errors.New("báº¡n khÃ´ng cÃ³ quyá»n xÃ³a danh má»¥c nÃ y")
	}

	return s.categoryRepo.Delete(categoryID)
}

