package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) GetMattermostUsers(c *gin.Context) {
	user, err := h.uc.GetMattermostUsers()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "ошибка предоставления списка кодов организаций", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) AddGLPICommentFromMattermostUser(c *gin.Context) {
	var commentFormMatt entity.MattermostInteractiveMessageRequestForm
	ip := ReadUserIP(c.Request)
	if !h.uc.MattermostIntegrationAllowed(ip) {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": "Комментарий не добавлен. Этот ip запрещен для запросов."})
		return
	}
	err := c.ShouldBindJSON(&commentFormMatt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": "Комментарий не добавлен. Ошибка определения параметров комментария"})
		return
	}

	comment, err := h.uc.AddGLPI_HRPTicketCommentFromMattermost(commentFormMatt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": err.Error()})
		return
	}
	var res = struct {
		Message string `json:"message,omitempty"`
		Props   any    `json:"props,omitempty"`
	}{
		Message: "> *" + comment + "*",
	}

	c.JSON(http.StatusOK, gin.H{"ephemeral_text": comment, "update": res})
}

func (h *Handler) DisableCalendarTaskNotification(c *gin.Context) {
	var commentFormMatt entity.MattermostInteractiveMessageRequestForm
	ip := ReadUserIP(c.Request)
	if !h.uc.MattermostIntegrationAllowed(ip) {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": "Оповещение не отменено. Этот ip запрещен для запросов."})
		return
	}
	err := c.ShouldBindJSON(&commentFormMatt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": "Оповещение не отменено. Ошибка определения параметров действия"})
		return
	}

	comment, err := h.uc.DisableSheduleTaskNotificationFromMattermost(commentFormMatt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ephemeral_text": err.Error()})
		return
	}
	var res = struct {
		Message string `json:"message,omitempty"`
		Props   any    `json:"props,omitempty"`
	}{
		Message: "> *" + comment + "*",
	}

	c.JSON(http.StatusOK, gin.H{"ephemeral_text": comment, "update": res})
}
