package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewADRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/ad")
	rg.GET("/users", h.TokenAuthMiddleware(), h.Users)
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.User)
	rg.GET("/computers", h.TokenAuthMiddleware(), h.Computers)
	rg.GET("/stats/counts", h.TokenAuthMiddleware(), h.GetAdCounts)

	return rg
}
