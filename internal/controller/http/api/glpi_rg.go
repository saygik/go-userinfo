package api

import "github.com/gin-gonic/gin"

func (h *Handler) NewGlpiRouterGroup() *gin.RouterGroup {
	rg := h.rg.Group("/glpi")
	rg.GET("/whoami", h.TokenAuthMiddleware(), h.GlpiCurrentUser)
	rg.GET("/user/:username", h.TokenAuthMiddleware(), h.GetGLPIUser)
	rg.GET("/nctickets", h.TokenAuthMiddleware(), h.GetTicketsNonClosed)
	rg.GET("/tickets/mygroups", h.TokenAuthMiddleware(), h.GetTicketsInMyGroups)
	rg.GET("/ticket/:id", h.TokenAuthMiddleware(), h.GetTicket)                                // * Заявка
	rg.GET("/ticket/solutions/:id", h.TokenAuthMiddleware(), h.GetGLPITicketSolutionTemplates) // * Шаблоны решений заявки
	rg.POST("/ticket/user/:id", h.TokenAuthMiddleware(), h.AddTicketUser)                      // * Добавление пользователя заявки
	rg.POST("/ticket", h.TokenAuthMiddleware(), h.AddTicket)                                   // * Комментарий заявки
	rg.POST("/comment/ticket/:id", h.TokenAuthMiddleware(), h.AddTicketComment)                // * Комментарий заявки
	rg.POST("/solution/ticket/:id", h.TokenAuthMiddleware(), h.AddTicketSolution)              // * Решение заявки
	rg.GET("/problem/:id", h.TokenAuthMiddleware(), h.GetProblem)                              // * Проблема
	rg.GET("/users", h.TokenAuthMiddleware(), h.GetGLPIUsers)
	rg.GET("/otkazes", h.TokenAuthMiddleware(), h.GetOtkazes)                                                // * Отказы
	rg.GET("/problems", h.TokenAuthMiddleware(), h.GetProblems)                                              // * Проблемы
	rg.GET("/statistics/tickets", h.TokenAuthMiddleware(), h.GetStatTickets)                                 // * Статистика по заявкам
	rg.GET("/statistics/failures", h.TokenAuthMiddleware(), h.GetStatFailures)                               // * Статистика по отказам
	rg.GET("/statistics/period-regions-month-days", h.TokenAuthMiddleware(), h.GetStatPeriodRegionDayCounts) // * Статистика по отказам
	rg.GET("/statistics/statsdays", h.TokenAuthMiddleware(), h.GetStatTicketsDays)                           // * Статистика по отказам
	rg.GET("/statistics/top10performers", h.TokenAuthMiddleware(), h.GetStatTop10Performers)
	rg.GET("/statistics/top10iniciators", h.TokenAuthMiddleware(), h.GetStatTop10Iniciators)
	rg.GET("/statistics/top10groups", h.TokenAuthMiddleware(), h.GetStatTop10Groups)
	rg.GET("/statistics/periodcounts", h.TokenAuthMiddleware(), h.GetStatPeriodTicketsCounts)
	rg.GET("/statistics/periodrequestypes", h.TokenAuthMiddleware(), h.GetStatPeriodRequestTypes)
	rg.GET("/statistics/regions", h.TokenAuthMiddleware(), h.GetStatRegions)
	rg.GET("/statistics/period-org-treemap", h.TokenAuthMiddleware(), h.GetStatPeriodOrgTreemap)

	rg.GET("/hrp", h.GetHRPTickets)

	return rg
}
