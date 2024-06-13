package api

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewOAuth2RouterGroup() *gin.RouterGroup {

	rg := h.rg.Group("/oauth")
	rg.GET("/login", h.GetLogin)
	rg.POST("/login", h.PostLogin)
	rg.GET("/consent", h.GetConsent)
	rg.POST("/consent", h.PostConsent)

	return rg
}
