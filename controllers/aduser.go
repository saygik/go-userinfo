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
func (ctrl ADUserController) AllDomains(c *gin.Context) {
	domains := aduserModel.AllDomains()
	c.JSON(http.StatusOK, gin.H{"data": domains})

}

//All ...
func (ctrl ADUserController) GroupUsers(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain group users", "error": "Empty domain name"})
		return
	}
	group := c.Param("group")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty group name"})
		return
	}
	users, err := aduserModel.GroupUsers(domain, group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get domain group users", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
