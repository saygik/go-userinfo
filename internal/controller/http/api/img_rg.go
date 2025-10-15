package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewImgRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/img")
	rg.GET("/ticket-status/:id", h.getImgTicketStatus)
	return rg
}
