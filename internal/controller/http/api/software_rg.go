package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewSoftwareRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/software")
	rg.GET("/one/:id", h.TokenAuthMiddleware(), h.GetSoftware)
	rg.POST("/one/:id", h.TokenAuthMiddleware(), h.AddUserToSoftware)
	rg.GET("/one/:id/users", h.TokenAuthMiddleware(), h.GetSoftwareUsers)
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.GetSoftwareUser)
	rg.POST("/user/:username", h.TokenAuthMiddleware(), h.AddSoftwareUser)
	rg.DELETE("/user/:id", h.TokenAuthMiddleware(), h.DelSoftwareUser)
	rg.GET("/all", h.TokenAuthMiddleware(), h.GetSoftwares)

	return rg
}
