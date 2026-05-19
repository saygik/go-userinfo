package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewIUTMRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/iutm")
	rg.GET("/wlist", func(c *gin.Context) {
		h.Wlist(c, "wlist")
	})
	rg.GET("/wlist2", func(c *gin.Context) {
		h.Wlist(c, "wlist2")
	})
	rg.GET("/wlistivc", func(c *gin.Context) {
		h.Wlist(c, "wlistivc")
	})
	rg.GET("/blist", func(c *gin.Context) {
		h.Wlist(c, "blist")
	})

	// rg.GET("/wlist", h.Wlist)
	// rg.GET("/wlist2", h.Wlist2)
	// rg.GET("/wlistivc", h.Wlistivc)
	// rg.GET("/blist", h.Blist)
	rg.GET("/alllists", h.AllLists)

	return rg
}
