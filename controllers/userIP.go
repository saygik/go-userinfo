package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/forms"
	"github.com/saygik/go-userinfo/models"
	"net/http"
	"strings"
)

//UserController ...
type UserIPController struct{}

var userIPModel = new(models.UserIPModel)
var DefaultDomain string

func (ctrl UserIPController) All(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		domain = DefaultDomain
	}
	results, err := userIPModel.All(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}

//
func (ctrl UserIPController) SetIp(c *gin.Context) {
	var userForm forms.UserActivityForm
	err := c.ShouldBindQuery(&userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": err.Error()})
		return

	}
	domain := strings.Split(userForm.User, "@")[1]

	if !ad.Domains[domain] {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set users ip", "error": "Non served domain"})
		return
	}

	userForm.Ip = ReadUserIP(c.Request)
	//	activity := query.Get("activity")
	if userForm.Activiy == "" {
		userForm.Activiy = "login"
	}
	var msgResponce string
	msgResponce, err = userIPModel.SetUserIp(userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set user ip", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msgResponce})
}

func (ctrl UserIPController) GetUserByName(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	userWithActivity, err := userIPModel.GetUserByName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userWithActivity})
	// ipAddr:=ReadUserIP(c.Request)
	//var hostNames []string
	// hostNames,_=ReadUserName(ipAddr)
	//c.JSON(http.StatusOK, gin.H{"data": ipAddr, "host_names": hostNames})

}
