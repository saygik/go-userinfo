package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewManualRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/manual")
	rg.GET("/orgcodes", h.TokenAuthMiddleware(), h.OrgCodes)

	return rg
}
