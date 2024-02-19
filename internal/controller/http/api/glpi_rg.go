package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewGlpiRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/glpi")
	rg.GET("/whoami", h.TokenAuthMiddleware(), h.GlpiCurrentUser)
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.GetGLPIUser)
	return rg
}
