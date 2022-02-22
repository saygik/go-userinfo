package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
	"net/http"
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
