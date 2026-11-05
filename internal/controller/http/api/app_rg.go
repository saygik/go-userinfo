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
	rg.POST("/admin/adcache/update", h.TokenAuthMiddleware(), h.ForceFillRedisCacheFromAD)
	rg.GET("/admin/adcache/status", h.TokenAuthMiddleware(), h.ForceFillRedisCacheFromADStatus)
	rg.POST("/admin/iutmcache/update", h.TokenAuthMiddleware(), h.ForceFillRedisCacheIUTM)
	rg.GET("/admin/overview", h.TokenAuthMiddleware(), h.AdminOverview)
	rg.DELETE("/admin/user/:username/rights", h.TokenAuthMiddleware(), h.AdminDeleteUserRights)
	rg.GET("/admin/audit", h.TokenAuthMiddleware(), h.AdminPermissionAudit)
	rg.GET("/admin/software-audit", h.TokenAuthMiddleware(), h.AdminUserSoftwareAudit)
	rg.GET("/admin/calendars", h.TokenAuthMiddleware(), h.AdminSchedules)
	rg.GET("/admin/calendars/meta", h.TokenAuthMiddleware(), h.AdminScheduleMeta)
	rg.POST("/admin/calendars", h.TokenAuthMiddleware(), h.AdminScheduleCreate)
	rg.PUT("/admin/calendars/:id", h.TokenAuthMiddleware(), h.AdminScheduleUpdate)
	rg.DELETE("/admin/calendars/:id", h.TokenAuthMiddleware(), h.AdminScheduleDelete)

	rg.GET("/ip", h.Ip)
	return rg
}
