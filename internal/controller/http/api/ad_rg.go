package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewADRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/ad")
	rg.GET("/users", h.TokenAuthMiddleware(), h.Users)
	rg.GET("/public/users", h.TokenAuthMiddleware(), h.PUsers)
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.User)
	rg.GET("/computers", h.TokenAuthMiddleware(), h.Computers)
	rg.GET("/stats/counts", h.TokenAuthMiddleware(), h.GetAdCounts)

	rg.GET("/activity/user/:username", h.TokenAuthMiddleware(), h.GetUserADActivity)
	rg.PUT("/user/avatar/:username", h.TokenAuthMiddleware(), h.UpdateUserAvatar)
	rg.PUT("/user/role/:username", h.TokenAuthMiddleware(), h.UpdateUserRole)
	rg.POST("/user/group/:username", h.TokenAuthMiddleware(), h.AddUserGroup)
	rg.DELETE("/user/group/:username", h.TokenAuthMiddleware(), h.DelUserGroup)

	return rg
}
