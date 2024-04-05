package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewMattermostRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/mattermost")
	rg.GET("/users", h.TokenAuthMiddleware(), h.GetMattermostUsers)

	return rg
}
