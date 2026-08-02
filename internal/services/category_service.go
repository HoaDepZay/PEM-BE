package services

import (
	"errors"
	"time"
	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
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

func (s *CategoryService) GetCategories(userID string) ([]models.CategoryGroup, error) {
	groups, err := s.categoryRepo.GetCategoryGroupsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Calculate Status and OverBudgetAmount dynamically for each Category
	now := time.Now()
	for i := range groups {
		for j := range groups[i].Categories {
			cat := &groups[i].Categories[j]
			var startDate, endDate time.Time

			switch cat.BudgetType {
			case "DAILY":
				startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				endDate = startDate.Add(24 * time.Hour)
			case "WEEKLY":
				offset := int(time.Monday - now.Weekday())
				if offset > 0 {
					offset = -6
				}
				startDate = time.Date(now.Year(), now.Month(), now.Day()+offset, 0, 0, 0, 0, now.Location())
				endDate = startDate.AddDate(0, 0, 7)
			case "MONTHLY":
				startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				endDate = startDate.AddDate(0, 1, 0)
			case "YEARLY":
				startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
				endDate = startDate.AddDate(1, 0, 0)
			default:
				cat.Status = "TỐT"
				continue
			}

			var totalSpent int64
			db.DB.Model(&models.Expense{}).
				Where("CategoryID = ? AND ExpenseDate >= ? AND ExpenseDate < ?", cat.CategoryID, startDate, endDate).
				Select("ISNULL(SUM(Amount), 0)").
				Scan(&totalSpent)

			if totalSpent > cat.BudgetAmount {
				cat.Status = "VƯỢT NGÂN SÁCH"
				cat.OverBudgetAmount = totalSpent - cat.BudgetAmount
			} else {
				cat.Status = "TỐT"
				cat.OverBudgetAmount = 0
			}
		}
	}

	return groups, nil
}

// ----------------- GROUPS -----------------
func (s *CategoryService) CreateGroup(userID, name, icon, color string) (*models.CategoryGroup, error) {
	if name == "" {
		return nil, errors.New("tên nhóm không được để trống")
	}

	uid, _ := uuid.Parse(userID)
	group := &models.CategoryGroup{
		GroupID: models.NewMSSQLUUID(),
		UserID:  (*models.MSSQLUUID)(&uid),
		Name:    name,
		Icon:    icon,
		Color:   color,
	}

	if err := s.categoryRepo.CreateGroup(group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *CategoryService) UpdateGroup(userID, groupID, name, icon, color string) (*models.CategoryGroup, error) {
	group, err := s.categoryRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("nhóm không tồn tại")
	}
	uid, _ := uuid.Parse(userID)
	if group.UserID == nil || uuid.UUID(*group.UserID) != uid {
		return nil, errors.New("bạn không có quyền sửa nhóm này")
	}

	if name != "" {
		group.Name = name
	}
	if icon != "" {
		group.Icon = icon
	}
	if color != "" {
		group.Color = color
	}

	if err := s.categoryRepo.UpdateGroup(group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *CategoryService) DeleteGroup(userID, groupID string) error {
	group, err := s.categoryRepo.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return errors.New("nhóm không tồn tại")
	}
	uid, _ := uuid.Parse(userID)
	if group.UserID == nil || uuid.UUID(*group.UserID) != uid {
		return errors.New("bạn không có quyền xóa nhóm này")
	}

	return s.categoryRepo.DeleteGroup(groupID)
}

// ----------------- CATEGORIES (Sub) -----------------
func (s *CategoryService) CreateCategory(userID, groupID, name, budgetType string, budgetAmount int64) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("tên danh mục không được để trống")
	}
	gID, err := uuid.Parse(groupID)
	if err != nil {
		return nil, errors.New("GroupID không hợp lệ")
	}

	// Xác minh quyền sở hữu group
	group, err := s.categoryRepo.GetGroupByID(groupID)
	if err != nil || group == nil {
		return nil, errors.New("nhóm không hợp lệ")
	}
	uid, _ := uuid.Parse(userID)
	if group.UserID != nil && uuid.UUID(*group.UserID) != uid {
		return nil, errors.New("bạn không có quyền thêm vào nhóm này")
	}

	category := &models.Category{
		CategoryID:   models.NewMSSQLUUID(),
		GroupID:      models.MSSQLUUID(gID),
		Name:         name,
		BudgetType:   budgetType,
		BudgetAmount: budgetAmount,
	}

	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return nil, err
	}
	
	// Cập nhật lại tổng ngân sách của Group
	s.categoryRepo.UpdateGroupBudget(groupID)
	
	return category, nil
}

func (s *CategoryService) UpdateCategory(userID, categoryID, name, budgetType string, budgetAmount int64) (*models.Category, error) {
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("danh mục không tồn tại")
	}
	uid, _ := uuid.Parse(userID)
	
	// Cần lấy Group của Category này để kiểm tra quyền
	group, err := s.categoryRepo.GetGroupByID(uuid.UUID(category.GroupID).String())
	if err != nil || group == nil || group.UserID == nil || uuid.UUID(*group.UserID) != uid {
		return nil, errors.New("bạn không có quyền sửa danh mục này")
	}

	if name != "" {
		category.Name = name
	}
	if budgetType != "" {
		category.BudgetType = budgetType
	}
	category.BudgetAmount = budgetAmount

	if err := s.categoryRepo.UpdateCategory(category); err != nil {
		return nil, err
	}
	
	// Cập nhật lại tổng ngân sách của Group
	s.categoryRepo.UpdateGroupBudget(uuid.UUID(category.GroupID).String())

	return category, nil
}

func (s *CategoryService) DeleteCategory(userID, categoryID string) error {
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("danh mục không tồn tại")
	}
	uid, _ := uuid.Parse(userID)

	group, err := s.categoryRepo.GetGroupByID(uuid.UUID(category.GroupID).String())
	if err != nil || group == nil || group.UserID == nil || uuid.UUID(*group.UserID) != uid {
		return errors.New("bạn không có quyền xóa danh mục này")
	}

	err = s.categoryRepo.DeleteCategory(categoryID)
	if err == nil {
		// Cập nhật lại tổng ngân sách của Group sau khi xóa
		s.categoryRepo.UpdateGroupBudget(uuid.UUID(category.GroupID).String())
	}
	return err
}
