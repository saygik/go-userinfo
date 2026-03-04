package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewAppRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/app")
	rg.GET("/whoami", h.TokenAuthMiddleware(), h.CurrentUser)
	// rg.GET("/userresources", h.TokenAuthMiddleware(), h.CurrentUserResources)
	// rg.GET("/resources", h.TokenAuthMiddleware(), h.AppResources)
	rg.GET("/roles", h.TokenAuthMiddleware(), h.AppRoles)
	rg.GET("/sections", h.TokenAuthMiddleware(), h.AppSections)
	rg.GET("/appdomains", h.TokenAuthMiddleware(), h.AppDomains)
	rg.GET("/domains", h.TokenAuthMiddleware(), h.DomainList)
	rg.GET("/computer/rms/:domain", h.ComputerRMS)
	rg.GET("/setip", h.SetIp)
	rg.POST("/localadmins/:computer", h.GetLocalAdmins)
	rg.POST("/computer-update-admins/:computer", h.TokenAuthMiddleware(), h.UpdateComputerLocalAdmins)
	rg.GET("/ip", h.Ip)
	return rg
}
