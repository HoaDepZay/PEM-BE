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

// GetCategories lấy danh sách category
// @Summary Lấy danh sách danh mục
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /categories [get]
func (c *CategoryController) GetCategories(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	categories, err := c.categoryService.GetCategories(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách danh mục"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory tạo danh mục mới
// @Summary Tạo danh mục mới
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Thông tin danh mục (name, icon, color)"
// @Success 201 {object} map[string]interface{}
// @Router /categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	var req struct {
		Name  string `json:"name" binding:"required"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		DailyBudget *int64 `json:"daily_budget"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	category, err := c.categoryService.CreateCategory(userIDStr, req.Name, req.Icon, req.Color, req.DailyBudget)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": category, "message": "Tạo danh mục thành công"})
}

// UpdateCategory sửa danh mục
// @Summary Sửa danh mục
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body map[string]string true "Thông tin danh mục cần sửa (name, icon, color)"
// @Success 200 {object} map[string]interface{}
// @Router /categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	userIDObj, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userIDObj.(uuid.UUID).String()

	categoryID := ctx.Param("id")

	var req struct {
		Name  string `json:"name"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		DailyBudget *int64 `json:"daily_budget"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	category, err := c.categoryService.UpdateCategory(userIDStr, categoryID, req.Name, req.Icon, req.Color, req.DailyBudget)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": category, "message": "Cập nhật danh mục thành công"})
}

// DeleteCategory xóa danh mục
// @Summary Xóa danh mục
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /categories/{id} [delete]
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Xóa danh mục thành công"})
}
