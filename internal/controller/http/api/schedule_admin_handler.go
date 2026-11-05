package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

// adminOnly проверяет, что запрос выполняется пользователем с ролью SysAdmin.
func (h *Handler) adminOnly(c *gin.Context) bool {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return false
	}
	if !perms.IsSysAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав для управления календарями"})
		return false
	}
	return true
}

// AdminSchedules отдаёт список календарей с выданными правами.
func (h *Handler) AdminSchedules(c *gin.Context) {
	if !h.adminOnly(c) {
		return
	}
	schedules, err := h.uc.GetAdminSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно получить список календарей",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": schedules})
}

// AdminScheduleMeta отдаёт справочники для формы календаря.
func (h *Handler) AdminScheduleMeta(c *gin.Context) {
	if !h.adminOnly(c) {
		return
	}
	groups, err := h.uc.GetScheduleGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно получить список групп GLPI",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"groups": groups}})
}

// AdminScheduleCreate создаёт календарь.
func (h *Handler) AdminScheduleCreate(c *gin.Context) {
	if !h.adminOnly(c) {
		return
	}
	var form entity.ScheduleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно создать календарь. Неправильные данные формы",
			"error":   err.Error(),
		})
		return
	}
	schedule, err := h.uc.CreateSchedule(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно создать календарь",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": schedule})
}

// AdminScheduleUpdate изменяет календарь и его права.
func (h *Handler) AdminScheduleUpdate(c *gin.Context) {
	if !h.adminOnly(c) {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неправильный идентификатор календаря"})
		return
	}
	var form entity.ScheduleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно изменить календарь. Неправильные данные формы",
			"error":   err.Error(),
		})
		return
	}
	if err := h.uc.UpdateSchedule(id, form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно изменить календарь",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Календарь изменён"})
}

// AdminScheduleDelete удаляет календарь вместе с задачами и правами.
func (h *Handler) AdminScheduleDelete(c *gin.Context) {
	if !h.adminOnly(c) {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неправильный идентификатор календаря"})
		return
	}
	if err := h.uc.DeleteSchedule(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно удалить календарь",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Календарь удалён"})
}
