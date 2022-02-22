package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
	"net/http"
)

//UserController ...
type ADUserController struct{}

var aduserModel = new(models.ADUserModel)

//All ...
func (ctrl ADUserController) All(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty domain name"})
		return
	}
	users, err := aduserModel.All(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get domain users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}
