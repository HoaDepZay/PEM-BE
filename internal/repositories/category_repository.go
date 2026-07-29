package repositories

import (
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"

	"gorm.io/gorm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

// GetCategoriesByUserID lấy danh sách các danh mục của 1 user và các danh mục hệ thống (UserID = NULL)
func (r *CategoryRepository) GetCategoriesByUserID(userID string) ([]models.Category, error) {
	var categories []models.Category
	// Lấy categories của user HOẶC hệ thống (UserID IS NULL)
	err := db.DB.Where("UserID = ?", userID).Or("UserID IS NULL").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return db.DB.Create(category).Error
}

func (r *CategoryRepository) Update(category *models.Category) error {
	return db.DB.Save(category).Error
}

func (r *CategoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := db.DB.First(&category, "CategoryID = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepository) Delete(id string) error {
	return db.DB.Delete(&models.Category{}, "CategoryID = ?", id).Error
}
