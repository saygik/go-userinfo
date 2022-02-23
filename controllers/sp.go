package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
	"net/http"
)

//UserController ...
type SPController struct{}

var spModel = new(models.SPModel)

//All ...
func (ctrl SPController) All(c *gin.Context) {

	presences, err := spModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get zals list from sharepoint server", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": presences})
}
