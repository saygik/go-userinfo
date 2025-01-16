package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewMattermostRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/mattermost")
	rg.GET("/users", h.TokenAuthMiddleware(), h.GetMattermostUsers)
	rg.POST("/glpi/comment", h.AddGLPICommentFromMattermostUser)
	rg.POST("/schedule/notification", h.DisableCalendarTaskNotification)

	return rg
}
