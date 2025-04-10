package api

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) NewOAuth2RouterGroup() *gin.RouterGroup {

	rg := h.rg.Group("/oauth-authentik")
	rg.GET("/login", h.LoginOauth2Authentik)
	rg.POST("/logout", h.LogoutOauth2Authentik)
	rg.GET("/token", h.ExchangeTokenAuthentik)
	rg.POST("/refresh", h.RefreshOauth2Authentik)
	return rg
}
