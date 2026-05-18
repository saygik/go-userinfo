package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewIUTMRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/iutm")
	rg.GET("/wlist", h.Wlist)
	rg.GET("/wlist2", h.Wlist2)
	rg.GET("/blist", h.Blist)
	rg.GET("/alllists", h.AllLists)

	return rg
}
