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
	newUsers := []models.MattermostUser{}
	for _, user := range users {
		if user.Roles != "system_user system_admin" && user.IsBot == false {
			user.Roles = ""
			if user.AuthService == "ldap" {
				user.AD = "it.rw"
			} else {
				user.AD = "-"
			}
			user.AuthService = ""
			//			user.IsBot = nil
			newUsers = append(newUsers, user)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": newUsers})
}
