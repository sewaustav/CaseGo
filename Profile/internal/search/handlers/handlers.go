package handlers

import (
	"net/http"

	_ "github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service service.SearchServiceInterface
}

func NewSearchHandler(service service.SearchServiceInterface) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

// GetProfilesHandler godoc
// @Summary Поиск профилей по фильтрам
// @Description Фильтрация профилей по профессии, возрасту, городу и полу
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profession_id query int false "ID профессии"
// @Param profession query string false "Название профессии"
// @Param min_age query int false "Мин. возраст"
// @Param max_age query int false "Макс. возраст"
// @Param city query string false "Город"
// @Param sex query int false "Пол (0 или 1)"
// @Param limit query int false "Лимит"
// @Param page query int false "Страница"
// @Success 200 {array} models.Profile
// @Failure 400 {object} map[string]string "Invalid search params"
// @Router /search [get]
func (h *SearchHandler) GetProfilesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.SearchDTO
	var helpers dto.SearchHelpersDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search params"})
		return
	}
	if err := c.ShouldBindQuery(&helpers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination params"})
		return
	}

	res, err := h.service.SearchProfileService(ctx, req, helpers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

// GetProfileFioHandler godoc
// @Summary Поиск профилей по ФИО
// @Description Поиск по имени, фамилии или отчеству
// @Tags search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param patronymic query string false "Отчество"
// @Param limit query int false "Лимит"
// @Param page query int false "Страница"
// @Success 200 {array} models.Profile
// @Router /search/fio [get]
func (h *SearchHandler) GetProfileFioHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.SearchByFIODTO
	var helpers dto.SearchHelpersDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search params"})
		return
	}
	if err := c.ShouldBindQuery(&helpers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination params"})
		return
	}

	res, err := h.service.SearchByFioService(ctx, req, helpers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
