package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) GlpiCurrentUser(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		user, err := h.uc.GetGlpiUser(userID)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}

}

func (h *Handler) GetGLPIUser(c *gin.Context) {
	userName := c.Param("username")
	userTechName := getUserID(c)
	user, err := h.uc.GetGlpiUserForTechnical(userName, userTechName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления профиля пользователя GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) GetTicketsNonClosed(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	tickets, err := h.uc.GetGLPITicketsNonClosed(user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка получения незакрытых заявок из GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
func (h *Handler) GetGLPIUsers(c *gin.Context) {

	users, err := h.uc.GetGLPIUsers()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "ошибка получения списка пользователей из GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *Handler) GetGLPITicketSolutionTemplates(c *gin.Context) {

	ticketId := c.Param("id")
	profiles, err := h.uc.GetGLPITicketSolutionTemplates(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка получения шаблонов решений из GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Шаблоны решений", "data": profiles})

}
func (h *Handler) GetTicket(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	ticketId := c.Param("id")
	ticket, err := h.uc.GetGLPITicket(ticketId, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка получения заявки из GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket finded", "data": ticket})

}
func (h *Handler) GetTicketReport(c *gin.Context) {
	user := getUserID(c)
	//user := "sb@brnv.rw"
	//user := ""
	if user == "" {
		//		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Для доступа сначала пройдите авторизацию", "data": entity.GLPI_Ticket_Report{}})
		return
	}
	ticketId := c.Param("id")
	ticket, err := h.uc.GetGLPITicketReport(ticketId, user)
	if err != nil {
		//c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка получения заявки из GLPI", "error": err.Error()})
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ticket})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "", "data": ticket})

}

func (h *Handler) AddTicketUser(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	var commentForm entity.GLPITicketUserForm

	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить пользователя GLPI: " + err.Error()})
		return
	}

	commentForm.User = user
	err = h.uc.AddTicketUser(commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка добавления пользователя заявки в GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь добавлен"})

}

func (h *Handler) AddTicket(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	var ticketForm entity.NewTicketForm

	err := c.ShouldBindJSON(&ticketForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить параметры добавляемой заявки: " + err.Error()})
		return
	}
	ticketForm.User = user
	id, err := h.uc.AddTicket(ticketForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка добавления заявки в GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": id, "message": "Комментарий добавлен"})
}

func (h *Handler) AddTicketComment(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	var commentForm entity.NewCommentForm

	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить добавляемую систему: " + err.Error()})
		return
	}

	commentForm.RequestTypesId = 11
	commentForm.ItemType = "Ticket"
	commentForm.User = user
	err = h.uc.AddTicketComment(commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка добавления комментария в GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Комментарий добавлен"})

}

func (h *Handler) GetOtkazes(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	otkazes, err := h.uc.GetGLPIOtkazes(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Отказы не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": otkazes})

}

func (h *Handler) GetStatTickets(c *gin.Context) {
	tickets, err := h.uc.GetStatTickets()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
func (h *Handler) GetStatFailures(c *gin.Context) {

	tickets, err := h.uc.GetStatFailures()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
func (h *Handler) GetStatPeriodRegionDayCounts(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	dat := strings.Split(enddate, "T")[0]
	date, error := time.Parse("2006-01-02", dat)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Неверная  дата конца периода"})
		return
	}
	day := date.Day()
	tickets, err := h.uc.GetStatPeriodRegionDayCounts(startdate, enddate, day)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (h *Handler) GetStatTicketsDays(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	tickets, err := h.uc.GetStatTicketsDays(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})

}

func (h *Handler) GetStatTop10Performers(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	stats, err := h.uc.GetStatTop10Performers(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователи не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetStatTop10Iniciators(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	stats, err := h.uc.GetStatTop10Iniciators(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователи не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetStatTop10Groups(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	stats, err := h.uc.GetStatTop10Groups(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователи не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetStatPeriodTicketsCounts(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	stats, err := h.uc.GetStatPeriodTicketsCounts(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Счетчики не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetStatPeriodRequestTypes(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	stats, err := h.uc.GetStatPeriodRequestTypes(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Счетчики не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handler) GetStatRegions(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	tickets, err := h.uc.GetStatRegions(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (h *Handler) GetStatPeriodOrgTreemap(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	tickets, err := h.uc.GetStatPeriodOrgTreemap(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (h *Handler) GetProblems(c *gin.Context) {
	startdate := c.Query("startdate")
	if startdate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала периода"})
		return
	}
	enddate := c.Query("enddate")
	if enddate == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата конца периода"})
		return
	}
	otkazes, err := h.uc.GetGLPIProblems(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Отказы не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": otkazes})

}
func (h *Handler) GetProblem(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	id := c.Param("id")
	ticket, err := h.uc.GetGLPIProblem(id, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "ошибка получения заявки из GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket finded", "data": ticket})

}
func (h *Handler) GetHRPTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "HRP bot activated"})
}

func (h *Handler) AddTicketSolution(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	var commentForm entity.NewCommentForm

	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "ошибка добавления решения в GLPI: " + err.Error()})
		return
	}

	commentForm.Status = 2
	commentForm.ItemType = "Ticket"
	commentForm.User = user
	err = h.uc.AddTicketSolution(commentForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "ошибка добавления решения в GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Комментарий добавлен"})

}

func (h *Handler) GetTicketsInMyGroups(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	tickets, err := h.uc.GetTicketsInExecutionGroups(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "ошибка получения заявок для ваших групп слежения в GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})

}
