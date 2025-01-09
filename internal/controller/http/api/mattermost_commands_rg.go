package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewMattermostCommandsRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/mattermost-commands")
	rg.POST("/glpi", h.MattermostGLPICommand)

	return rg
}
