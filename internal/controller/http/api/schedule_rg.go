package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewScheduleRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/schedule")
	rg.GET("/one/:id", h.GetSchedule)
	rg.GET("/task/:id", h.GetScheduleTasks)
	rg.POST("/task", h.AddScheduleTask)
	rg.DELETE("/task/:id", h.DelScheduleTask)
	rg.PUT("/task/:id", h.UpdateScheduleTask)
	return rg
}
