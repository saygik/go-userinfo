package controllers

import (
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
