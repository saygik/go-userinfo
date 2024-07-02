package api

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewOAuth2RouterGroup() *gin.RouterGroup {

	rg := h.rg.Group("/oauth")
	rg.GET("/login", h.LoginOauth2)
	rg.POST("/logout", h.TokenAuthMiddleware(), h.LogoutOauth2)
	rg.GET("/token", h.ExchangeToken)
	rg.POST("/refresh", h.RefreshOauth2)
	return rg
}
