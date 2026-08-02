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

// Lấy danh sách nhóm danh mục (kèm danh mục con)
func (r *CategoryRepository) GetCategoryGroupsByUserID(userID string) ([]models.CategoryGroup, error) {
	var groups []models.CategoryGroup
	err := db.DB.Preload("Categories").Where("UserID = ? OR UserID IS NULL", userID).Find(&groups).Error
	return groups, err
}

// Lấy Group theo ID
func (r *CategoryRepository) GetGroupByID(id string) (*models.CategoryGroup, error) {
	var group models.CategoryGroup
	err := db.DB.First(&group, "GroupID = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (r *CategoryRepository) CreateGroup(group *models.CategoryGroup) error {
	return db.DB.Create(group).Error
}

func (r *CategoryRepository) UpdateGroup(group *models.CategoryGroup) error {
	return db.DB.Save(group).Error
}

func (r *CategoryRepository) DeleteGroup(id string) error {
	return db.DB.Delete(&models.CategoryGroup{}, "GroupID = ?", id).Error
}

// Các hàm cho Category con
func (r *CategoryRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	err := db.DB.First(&category, "CategoryID = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	return db.DB.Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	return db.DB.Save(category).Error
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	return db.DB.Delete(&models.Category{}, "CategoryID = ?", id).Error
}

// Cập nhật lại tổng ngân sách của một Group bằng cách tính tổng BudgetAmount của các Category con
func (r *CategoryRepository) UpdateGroupBudget(groupID string) error {
	return db.DB.Exec(`
		UPDATE CategoryGroups 
		SET TotalBudget = ISNULL((SELECT SUM(BudgetAmount) FROM Categories WHERE GroupID = ?), 0)
		WHERE GroupID = ?`, 
		groupID, groupID,
	).Error
}
