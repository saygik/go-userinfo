package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewAppRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/app")
	rg.GET("/whoami", h.TokenAuthMiddleware(), h.CurrentUser)
	rg.GET("/resources", h.TokenAuthMiddleware(), h.CurrentUserResources)
	rg.GET("/domains", h.TokenAuthMiddleware(), h.DomainList)

	return rg
}
