package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewSoftwareRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/software")
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.GetSoftwareUser)

	return rg
}
