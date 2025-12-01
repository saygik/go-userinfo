package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewIUTMRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/iutm")
	rg.GET("/wlist", h.Wlist)

	return rg
}
