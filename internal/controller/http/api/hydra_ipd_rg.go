package api

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewHydraIDPRouterGroup() *gin.RouterGroup {

	rg := h.rg.Group("/idp")
	rg.GET("/login", h.GetLogin)
	rg.POST("/login", h.PostLogin)
	rg.GET("/logout", h.GetLogout)
	rg.POST("/logout", h.PostLogout)
	rg.GET("/consent", h.GetConsent)
	rg.POST("/consent", h.PostConsent)
	return rg
}
