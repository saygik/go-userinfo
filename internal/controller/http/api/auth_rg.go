package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewAuthRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/auth")
	rg.POST("/login", h.Login)
	rg.GET("/logout", h.Logout)
	rg.POST("/refresh", h.Refresh)
	return rg
}
