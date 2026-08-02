package controllers

import (
	"net/http"
	"visualfinance/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: services.NewCategoryService(),
	}
}

// GetCategories lấy danh sách category groups kèm categories con
func (c *CategoryController) GetCategories(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	groups, err := c.categoryService.GetCategories(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách danh mục"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": groups})
}

// ------------------- GROUPS -------------------

// CreateGroup
func (c *CategoryController) CreateGroup(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	var req struct {
		Name  string `json:"name" binding:"required"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	group, err := c.categoryService.CreateGroup(userIDStr, req.Name, req.Icon, req.Color)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": group, "message": "Tạo nhóm danh mục thành công"})
}

// UpdateGroup
func (c *CategoryController) UpdateGroup(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()
	groupID := ctx.Param("id")

	var req struct {
		Name  string `json:"name"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	group, err := c.categoryService.UpdateGroup(userIDStr, groupID, req.Name, req.Icon, req.Color)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": group, "message": "Cập nhật nhóm danh mục thành công"})
}

// DeleteGroup
func (c *CategoryController) DeleteGroup(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()
	groupID := ctx.Param("id")

	err := c.categoryService.DeleteGroup(userIDStr, groupID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Xóa nhóm danh mục thành công"})
}

// ------------------- SUB-CATEGORIES -------------------

// CreateCategory
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	var req struct {
		GroupID      string `json:"group_id" binding:"required"`
		Name         string `json:"name" binding:"required"`
		BudgetType   string `json:"budget_type" binding:"required"`
		BudgetAmount int64  `json:"budget_amount"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	category, err := c.categoryService.CreateCategory(userIDStr, req.GroupID, req.Name, req.BudgetType, req.BudgetAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": category, "message": "Tạo danh mục con thành công"})
}

// UpdateCategory
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()
	categoryID := ctx.Param("id")

	var req struct {
		Name         string `json:"name"`
		BudgetType   string `json:"budget_type"`
		BudgetAmount int64  `json:"budget_amount"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	category, err := c.categoryService.UpdateCategory(userIDStr, categoryID, req.Name, req.BudgetType, req.BudgetAmount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": category, "message": "Cập nhật danh mục con thành công"})
}

// DeleteCategory
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()
	categoryID := ctx.Param("id")

	err := c.categoryService.DeleteCategory(userIDStr, categoryID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Xóa danh mục con thành công"})
}
