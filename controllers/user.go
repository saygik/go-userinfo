package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
	"net/http"
)

//UserController ...
type UserController struct{}

var userModel = new(models.UserModel)

func (ctrl UserController) All(c *gin.Context) {
	results, err := userModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
	// ipAddr:=ReadUserIP(c.Request)
	//var hostNames []string
	// hostNames,_=ReadUserName(ipAddr)
	//c.JSON(http.StatusOK, gin.H{"data": ipAddr, "host_names": hostNames})

}
