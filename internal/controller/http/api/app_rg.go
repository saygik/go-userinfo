package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewAppRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/app")
	rg.GET("/whoami", h.TokenAuthMiddleware(), h.CurrentUser)
	rg.GET("/userresources", h.TokenAuthMiddleware(), h.CurrentUserResources)
	rg.GET("/resources", h.TokenAuthMiddleware(), h.AppResources)
	rg.GET("/roles", h.TokenAuthMiddleware(), h.AppRoles)
	rg.GET("/groups", h.TokenAuthMiddleware(), h.AppGroups)
	rg.GET("/domains", h.TokenAuthMiddleware(), h.DomainList)
	rg.GET("/setip", h.SetIp)
	return rg
}
