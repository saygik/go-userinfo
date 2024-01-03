package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type MattermostController struct{}

var mattModel = new(models.MattermostModel)

// All ...
func (ctrl MattermostController) GetAll(c *gin.Context) {

	users, err := mattModel.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users list from Mattermost server", "error": err.Error()})
		return
	}
	softUsers, err := userIPModel.GetSoftwareUsers(833)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно получить список пользователей Mattermost, зарегистрированных в системе", "error_message": err.Error()})
		return
	}

	for i, user := range users {
		if user.Roles != "system_user system_admin" && user.IsBot == false {
			users[i].Roles = ""

			if user.AuthService == "ldap" {
				users[i].AD = "it.rw"
			} else {
				users[i].AD = "-"
			}
			users[i].Registred = false
			for _, softUser := range softUsers {
				if softUser.Login == user.Name {
					users[i].Registred = true
				}
			}
			users[i].AuthService = ""
			//			user.IsBot = nil

		}
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
