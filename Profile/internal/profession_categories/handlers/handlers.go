package handlers

import (
	"net/http"
	"strconv"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/service"
	userModel "github.com/YoungFlores/Case_Go/Profile/internal/profile/models"

	"github.com/gin-gonic/gin"
)

type ProfessionCategoryHandler struct {
	service *categoryService.ProfessionCategoryService
}

func NewProfessionCategoryHandler(service *categoryService.ProfessionCategoryService) *ProfessionCategoryHandler {
	return &ProfessionCategoryHandler{
		service: service,
	}
}


func (h *ProfessionCategoryHandler) GetRole(c *gin.Context) (userModel.UserRole, bool) {
	userRole, exist := c.Get("role")
	if !exist {
		return 0, false
	}

	userRoleInt, ok := userRole.(userModel.UserRole)
	if !ok {
		return 0, false
	}

	return userRoleInt, true

}

// CreateCategoryHandler godoc
// @Summary Создать категорию профессий
// @Description Создает новую категорию (только для админов)
// @Tags profession-category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CategoryDTO true "Данные категории"
// @Success 201 {object} models.CategoryDTO
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /categories [post]
func (h *ProfessionCategoryHandler) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userRole, exist := h.GetRole(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if userRole != userModel.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var dto models.CategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.service.CreateCategoryService(ctx, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)

}

// GetCategoriesHandler godoc
// @Summary Получить все категории
// @Description Возвращает полный список категорий профессий
// @Tags profession-category
// @Produce json
// @Success 200 {array} models.CategoryDTO
// @Router /categories [get]
func (h *ProfessionCategoryHandler) GetCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	category, err := h.service.GetCategoriesService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryByParentIDHandler godoc
// @Summary Получить подкатегории по ParentID
// @Description Возвращает список категорий, привязанных к конкретному родителю
// @Tags profession-category
// @Produce json
// @Param parentID path int true "Parent ID"
// @Success 200 {array} models.CategoryDTO
// @Failure 400 {object} map[string]string "Invalid parentID"
// @Router /categories/parent/{parentID} [get]
func (h *ProfessionCategoryHandler) GetCategoryByParentIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	parentID, err := strconv.Atoi(c.Param("parentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parentID"})
		return
	}

	categories, err := h.service.GetCategoriesByParentService(ctx, int16(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)

}

// GetCategoryByIDHandler godoc
// @Summary Получить категорию по ID
// @Description Возвращает одну категорию по её уникальному идентификатору
// @Tags profession-category
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.CategoryDTO
// @Failure 400 {object} map[string]string "Invalid id"
// @Router /categories/{id} [get]
func (h *ProfessionCategoryHandler) GetCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	category, err := h.service.GetCategoryByIDService(ctx, int16(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)

}
