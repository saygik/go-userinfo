package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
	"net/http"
	"strconv"
)

//UserController ...
type SkypeController struct{}

var skypeModel = new(models.SkypeModel)

//All ...
func (ctrl SkypeController) AllPresences(c *gin.Context) {

	presences, err := skypeModel.AllPresences()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get user presences from skype server", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": presences})
}

//OnePresence ...
func (ctrl SkypeController) OnePresence(c *gin.Context) {

	user := c.Param("user")
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	userPresence, err := skypeModel.OnePresence(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Article not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userPresence})
}

//AllActiveConferences
func (ctrl SkypeController) AllActiveConferences(c *gin.Context) {

	data, err := skypeModel.AllActiveConferences()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get active conferences from skype server", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

//ConferencePresence
func (ctrl SkypeController) ConferencePresence(c *gin.Context) {
	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {
		conferencePresence, err := skypeModel.ConferencePresence(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Could not get users for conference from skype server", "error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"data": conferencePresence})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
	}
}
