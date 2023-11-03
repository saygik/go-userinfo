package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/forms"
	"github.com/saygik/go-userinfo/models"
)

type GLPIController struct{}

var GLPIModel = new(models.GLPIModel)

func (ctrl GLPIController) GetUsers(c *gin.Context) {

	users, err := GLPIModel.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список пользователей GLPI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})

}

// ************************* Всего отказов **********************************//
func (ctrl GLPIController) GetStatOtkazSum(c *gin.Context) {
	sum, _ := GLPIModel.GetStatOtkazSum()

	c.JSON(http.StatusOK, gin.H{"data": sum})
}

func (ctrl GLPIController) GetOtkazes(c *gin.Context) {
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
	otkazes, err := GLPIModel.GetOtkazes(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Отказы не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": otkazes})
}

func (ctrl GLPIController) CurrentUserGLPI(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	glpiUser, err := GLPIModel.GetUserByName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователь не найден в системе GLPI", "error": err.Error()})
		return
	}
	glpiUserProfiles, err := GLPIModel.GetUserProfiles(glpiUser.Id)
	if err == nil {
		glpiUser.Profiles = glpiUserProfiles
	}
	glpiUserGroups, err := GLPIModel.GetUserGroups(glpiUser.Id)
	if err == nil {
		glpiUser.Groups = glpiUserGroups
	}

	c.JSON(http.StatusOK, gin.H{"data": glpiUser})

}
func (ctrl GLPIController) GetUserByName(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}

	userID := getUserID(c)
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	userRequesterInGLPI, err := GLPIModel.GetUserByName(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Вы не найдены в системе GLPI"})
		return
	}
	glpiUserRequesterProfiles, err := GLPIModel.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Ваш профиль не найден в системе GLPI"})
		return
	}

	glpiUser, err := GLPIModel.GetUserByName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователь не найден в системе GLPI", "error": err.Error()})
		return
	}
	glpiUserProfiles, err := GLPIModel.GetUserProfiles(glpiUser.Id)
	if err == nil {
		glpiUser.Profiles = glpiUserProfiles
	}
	if !isTechnicalAdminOfUser(glpiUser, userRequesterInGLPI) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на этого пользователя в системе GLPI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": glpiUser})

}
func (ctrl GLPIController) GetSoftwares(c *gin.Context) {

	softwares, err := GLPIModel.GetSoftwares()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get softwares info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": softwares})

}
func (ctrl GLPIController) GetSoftware(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить имя пользователя в запросе"})
		return
	}

	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Неправильный идентификатор программного обеспечения"})
		return
	}
	softwares, err := GLPIModel.GetSoftware(id_int)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get software info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": softwares})

}
func (ctrl GLPIController) GetSoftwareUsers(c *gin.Context) {
	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {

		users, err := userIPModel.GetSoftwareUsers(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список пользователей программного обеспечения", "error_message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": users})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Неправильный номер программного обеспечения"})
	}

}
func (ctrl GLPIController) GetUserSoftwares(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить имя пользователя в запросе"})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Неправильное имя пользователя. Имя должно содержать домен в формате: user@domain"})
		return
	}
	// Список систем
	softwares, err := GLPIModel.GetSoftwares()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список систем из GLPI"})
		return
	}
	// список администраторов систем
	softwares = softwares
	userSoftwares, err := userIPModel.GetUserSoftwares(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список систем пользователя"})
		return
	}
	filteredSoft := []models.Software{}
	for _, soft := range softwares {
		for _, idsoft := range userSoftwares {
			if soft.Id == idsoft {
				filteredSoft = append(filteredSoft, soft)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": filteredSoft})

}

func (ctrl GLPIController) DelOneUserSoftware(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить имя пользователя в запросе"})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Неправильное имя пользователя. Имя должно содержать домен в формате: user@domain"})
		return
	}

	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {

		ra, err := userIPModel.DelOneUserSoftware(user, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Ошибка удаления система из списка пользователя: " + err.Error()})
			return
		}
		if ra == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Система отсутствует в списке пользователя"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cистема удалена из списка пользователя"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Неправильный номер системы"})
	}
}

func (ctrl GLPIController) AddOneUserSoftware(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить имя пользователя в запросе"})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Неправильное имя пользователя. Имя должно содержать домен в формате: user@domain"})
		return
	}
	var softwareForm forms.SoftwareForm
	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&softwareForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить номер системы: " + err.Error()})
		return
	}
	softwareForm.User = user
	rowsAffected, err := userIPModel.AddOneUserSoftware(softwareForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Ошибка добавления системы: " + err.Error()})
		return
	}
	if rowsAffected < 1 {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Система не добавлена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Система добавлена"})

}

func (ctrl GLPIController) GetStatTickets(c *gin.Context) {

	tickets, err := GLPIModel.GetStatTickets()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})

}
func (ctrl GLPIController) GetStatFailures(c *gin.Context) {

	tickets, err := GLPIModel.GetStatFailures()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (ctrl GLPIController) GetStatRegions(c *gin.Context) {
	date := c.Query("startdate")
	if date == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указана дата начала сбора статистики"})
		return
	}
	tickets, err := GLPIModel.GetStatRegions(date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (ctrl GLPIController) GetTicketsNonClosed(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	userRequesterInGLPI, err := GLPIModel.GetUserByName(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Вы не найдены в системе GLPI"})
		return
	}
	glpiUserRequesterProfiles, err := GLPIModel.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		if len(glpiUserRequesterProfiles) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Ваш профиль не найден в системе GLPI"})
			return
		}
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Ваш профиль не найден в системе GLPI"})
		return
	}
	glpiUserGroups, err := GLPIModel.GetUserGroups(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Groups = glpiUserGroups
	} else {
		userRequesterInGLPI.Groups = []models.IdName{}
	}
	tickets, err := GLPIModel.GetTicketsNonClosed()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get not closed tickets from GLPI", "error": err.Error()})
		return
	}
	var ticketOrgs []int64
	var ticketUsers []models.IdNameType
	var ticketGroups []models.IdNameType
	ticketsAllowed := []models.Ticket{}
	for _, ticket := range tickets {
		ticket.MyTicket = 0
		if err := json.Unmarshal([]byte(ticket.Users), &ticketUsers); err == nil {
			if containsInt64InIdNameTypeArray(ticketUsers, userRequesterInGLPI.Id) {
				ticket.MyTicket = 1
			}
		}
		if err := json.Unmarshal([]byte(ticket.UserGroups), &ticketGroups); err == nil {
			if containsIDNameInIdNameTypeArray(ticketGroups, userRequesterInGLPI.Groups) {
				ticket.MyTicket = 2
			}
		}
		for _, tp := range userRequesterInGLPI.Profiles {
			if tp.Id == 6 {
				if ticket.MyTicket > 0 {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break
				}
			}
			if tp.Id != 3 && tp.Id != 4 && tp.Id != 5 && tp.Id != 15 && tp.Id != 7 {
				continue
			}
			if tp.Recursive {
				if err := json.Unmarshal([]byte(ticket.Orgs), &ticketOrgs); err != nil {
					continue
				}
				if ticket.Eid == tp.Eid {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break

				}
				if containsInt64(ticketOrgs, tp.Eid) {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break

				}
			} else {
				if ticket.Eid == tp.Eid {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break
				}
			}
			// if tp.Recursive {
			// }
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": ticketsAllowed})
}

func (ctrl GLPIController) GetTicket(c *gin.Context) {
	ticketId := c.Param("id")
	if ticketId == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get ticket id from request", "error": ""})
		return
	}
	ticket, err := GLPIModel.GetTicket(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get ticket from GLPI", "error": err.Error()})
		return
	}
	works, _ := GLPIModel.TicketWorks(ticketId)
	ticket.Works = works
	c.JSON(http.StatusOK, gin.H{"message": "Ticket finded", "data": ticket})

}
