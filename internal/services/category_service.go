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

func (s *CategoryService) CreateCategory(userID, name, icon, color string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("tên danh mục không được để trống")
	}

	uid, _ := uuid.Parse(userID)

	category := &models.Category{
		CategoryID: uuid.New(),
		UserID: &uid,
		Name:   name,
		Icon:   icon,
		Color:  color,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(userID, categoryID, name, icon, color string) (*models.Category, error) {
	// Kiểm tra xem category có tồn tại và thuộc về user không
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("danh mục không tồn tại")
	}
	uid, _ := uuid.Parse(userID)
	if category.UserID == nil || *category.UserID != uid {
		return nil, errors.New("bạn không có quyền sửa danh mục này")
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
		return errors.New("danh mục không tồn tại")
	}
	uid, _ := uuid.Parse(userID)
	if category.UserID == nil || *category.UserID != uid {
		return errors.New("bạn không có quyền xóa danh mục này")
	}

	return s.categoryRepo.Delete(categoryID)
}
