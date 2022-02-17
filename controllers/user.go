package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/forms"
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

//
func (ctrl UserController) SetIp(c *gin.Context) {
	var userForm forms.UserActivityForm
	err := c.ShouldBindQuery(&userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": err.Error()})
		return

	}
	//query := c.Request.URL.Query()
	//user := query.Get("user")
	//
	//if user == "" {
	//	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
	//	return
	//}
	//if !isEmailValid(user) {
	//	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
	//	return
	//}
	userForm.Ip = ReadUserIP(c.Request)
	//	activity := query.Get("activity")
	if userForm.Activiy == "" {
		userForm.Activiy = "login"
	}
	err = userModel.SetUserIp(userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set user ip", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userForm.User, "ip": userForm.Ip, "activity": userForm.Activiy})
}
