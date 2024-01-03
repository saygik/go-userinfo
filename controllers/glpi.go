package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/db"
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
	for i, user := range users {
		adUser := GetRedisUser(user.Name)
		users[i].ADProfile = adUser
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
func (ctrl GLPIController) GetStatTop10Performers(c *gin.Context) {
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
	stats, err := GLPIModel.GetStatTop10Performers(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователи не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
func (ctrl GLPIController) GetStatPeriodTicketsCounts(c *gin.Context) {
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
	stats, err := GLPIModel.GetStatPeriodTicketsCounts(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Счетчики не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
func (ctrl GLPIController) GetStatPeriodRequestTypes(c *gin.Context) {
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
	stats, err := GLPIModel.GetStatPeriodRequestTypes(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Счетчики не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (ctrl GLPIController) GetStatTop10Iniciators(c *gin.Context) {
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
	stats, err := GLPIModel.GetStatTop10Iniciators(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Пользователи не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
func (ctrl GLPIController) GetStatTop10Groups(c *gin.Context) {
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
	stats, err := GLPIModel.GetStatTop10Groups(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Группы не найдены в системе GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
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
	admins, _ := GLPIModel.GetSoftwaresAdmins()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": softwares})
	}
	softAdmins := []map[string]interface{}{}
	for i, soft := range softwares {
		for _, admin := range admins {
			if soft.Groups_id_tech == admin.Id {
				adUser := GetRedisUser(admin.Name)
				softAdmins = append(softAdmins, adUser)
			}
		}
		if len(softAdmins) > 0 {
			softwares[i].Admins = softAdmins
			softAdmins = []map[string]interface{}{}
		} else {
			soft.Admins = []map[string]interface{}{}
		}

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
	software, err := GLPIModel.GetSoftware(id_int)

	//redisClient.HSet(ctx, "allusers", user["userPrincipalName"], jsonUser).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get software info from GLPI", "error": err.Error()})
		return
	}
	admins, _ := GLPIModel.GetSoftwaresAdmins()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": software})
	}
	softAdmins := []map[string]interface{}{}

	for _, admin := range admins {
		if software.Groups_id_tech == admin.Id {
			adUser := GetRedisUser(admin.Name)
			softAdmins = append(softAdmins, adUser)
		}
	}
	if len(softAdmins) > 0 {
		software.Admins = softAdmins
		softAdmins = nil
	} else {
		software.Admins = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"data": software})

}
func GetRedisUser(name string) (user map[string]interface{}) {
	var userProperties map[string]interface{}
	redisClient := db.GetRedis()
	var ctx = context.Background()
	redisADUser, err := redisClient.HGet(ctx, "allusers", name).Result()
	if err == nil {
		json.Unmarshal([]byte(redisADUser), &userProperties)
		delete(userProperties, "ip")
		delete(userProperties, "pwdLastSet")
		delete(userProperties, "proxyAddresses")
		delete(userProperties, "passwordDontExpire")
		delete(userProperties, "passwordCantChange")
		delete(userProperties, "distinguishedName")
		delete(userProperties, "userAccountControl")
		delete(userProperties, "memberOf")
		delete(userProperties, "employeeNumber")
		delete(userProperties, "presence")
		delete(userProperties, "otherTelephone")
		userProperties["name"] = name
		userProperties["findedInAD"] = true
		return userProperties
	} else {
		userProperties = map[string]interface{}{}
		userProperties["name"] = name
		userProperties["findedInAD"] = false

		return userProperties
	}
}
func (ctrl GLPIController) GetSoftwareUsers(c *gin.Context) {
	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {

		users, err := userIPModel.GetSoftwareUsers(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список пользователей программного обеспечения", "error_message": err.Error()})
			return
		}
		softUsers := []map[string]interface{}{}
		for _, user := range users {
			adUser := GetRedisUser(user.Name)
			adUser["login"] = user.Login
			adUser["comment"] = user.Comment
			adUser["fio"] = user.Fio
			adUser["external"] = user.External
			softUsers = append(softUsers, adUser)
		}

		c.JSON(http.StatusOK, gin.H{"data": softUsers})
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
	admins, _ := GLPIModel.GetSoftwaresAdmins()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": softwares})
	}
	softAdmins := []map[string]interface{}{}
	for i, soft := range softwares {
		for _, admin := range admins {
			if soft.Groups_id_tech == admin.Id {
				adUser := GetRedisUser(admin.Name)
				softAdmins = append(softAdmins, adUser)
			}
		}
		if len(softAdmins) > 0 {
			softwares[i].Admins = softAdmins
			softAdmins = []map[string]interface{}{}
		} else {
			soft.Admins = []map[string]interface{}{}
		}

	}

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

func (ctrl GLPIController) AddOneSoftwareUser(c *gin.Context) {
	id := c.Param("software")
	var idd int64
	if idx, err := strconv.ParseInt(id, 10, 64); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить номер системы"})
		return
	} else {
		idd = idx
	}
	var softwareForm forms.SoftwareUsersForm
	err := c.ShouldBindJSON(&softwareForm)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить пользователя системы: " + err.Error()})
		return
	}
	if softwareForm.User == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить имя пользователя в запросе"})
		return
	}
	if !isEmailValid(softwareForm.User) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Неправильное имя пользователя. Имя должно содержать домен в формате: user@domain"})
		return
	}

	softwareForm.Id = idd
	rowsAffected, err := userIPModel.AddOneSoftwareUser(softwareForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Ошибка добавления пользователя системы: " + err.Error()})
		return
	}
	if rowsAffected < 1 {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Пользователь системы не добавлен"})
		return
	}
	userProperties := GetRedisUser(softwareForm.User)
	delete(userProperties, "ip")
	delete(userProperties, "pwdLastSet")
	delete(userProperties, "proxyAddresses")
	delete(userProperties, "passwordDontExpire")
	delete(userProperties, "passwordCantChange")
	delete(userProperties, "distinguishedName")
	delete(userProperties, "userAccountControl")
	delete(userProperties, "memberOf")
	delete(userProperties, "employeeNumber")
	delete(userProperties, "presence")
	delete(userProperties, "url")
	delete(userProperties, "otherTelephone")
	userProperties["name"] = softwareForm.User
	userProperties["login"] = softwareForm.Login
	userProperties["comment"] = softwareForm.Comment
	userProperties["fio"] = softwareForm.Fio
	userProperties["external"] = softwareForm.External

	_, ok := userProperties["userPrincipalName"]
	if ok {
		userProperties["findedInAD"] = true
	} else {
		userProperties["findedInAD"] = false
	}
	c.JSON(http.StatusOK, gin.H{"data": userProperties})

}
func (ctrl GLPIController) GetStatTickets(c *gin.Context) {

	tickets, err := GLPIModel.GetStatTickets()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})

}
func (ctrl GLPIController) GetStatTicketsDays(c *gin.Context) {
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
	tickets, err := GLPIModel.GetStatTicketsDays(startdate, enddate)
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
	tickets, err := GLPIModel.GetStatRegions(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (ctrl GLPIController) GetStatPeriodOrgTreemap(c *gin.Context) {
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
	tickets, err := GLPIModel.GetStatPeriodOrgTreemap(startdate, enddate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get tickets statistics info from GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (ctrl GLPIController) GetStatPeriodRegionDayCounts(c *gin.Context) {
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
	tickets, err := GLPIModel.GetStatPeriodRegionDayCounts(startdate, enddate, day)
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
