package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
)

type GLPIController struct{}

var GLPIModel = new(models.GLPIModel)

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
	glpiUser, err := GLPIModel.GetUserByName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users info from GLPI", "error": err.Error()})
		return
	}
	glpiUserProfiles, err := GLPIModel.GetUserProfiles(glpiUser.Id)
	if err == nil {
		glpiUser.Profiles = glpiUserProfiles
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
