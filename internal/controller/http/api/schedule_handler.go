package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) GetSchedule(c *gin.Context) {
	user := ""
	userI, exists := c.Get("user")
	if exists {
		user = userI.(string)
	}
	_ = user
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указан ID расписания"})
		return
	}
	results, err := h.uc.GetSchedule(id, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить календарь", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}

func (h *Handler) GetScheduleTasks(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указан ID расписания"})
		return
	}
	results, err := h.uc.GetScheduleTasks(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить задачи календаря", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}

func (h *Handler) DelScheduleTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно обновить запись календаря. Неправильный id календаря", "error": err.Error()})
		return
	}

	err = h.uc.DelScheduleTask(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Запись календаря не удалена", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Запись удалена"})
}
func (h *Handler) AddScheduleTask(c *gin.Context) {

	var scheduleForm entity.ScheduleTask
	err := c.ShouldBindJSON(&scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно создать запись календаря. Неправильные данные формы", "error": err.Error()})
		return
	}

	msgResponce, err := h.uc.AddScheduleTask(scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно создать запись календаря", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": msgResponce})

}

func (h *Handler) UpdateScheduleTask(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно обновить запись календаря. Неправильный id календаря", "error": err.Error()})
		return
	}
	var scheduleForm entity.ScheduleTask
	err = c.ShouldBindJSON(&scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно обновить запись календаря. Неправильные данные формы", "error": err.Error()})
		return
	}
	scheduleForm.Id = id
	err = h.uc.UpdateScheduleTask(scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Запись календаря не обновлена", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Событие календаря изменено"})

}
