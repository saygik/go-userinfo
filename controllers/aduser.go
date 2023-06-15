package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type ADUserController struct{}

var aduserModel = new(models.ADUserModel)

func (ctrl ADUserController) GetAllDomainsUsers(clearDomains bool) {
	fmt.Println("Getting all Users...")
	if clearDomains {
		aduserModel.ClearAllDomainsUsers()
	}
	aduserModel.GetAllDomainsUsers()
}

// getUserID ...
func getUserID(c *gin.Context) (userID string) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	_, isExist := c.Get("user")
	if isExist {
		return c.MustGet("user").(models.UserInRedis).Login
	}
	return ""
}

// getUser ...
func getUser(c *gin.Context) (user models.UserInRedisOpenID) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	return c.MustGet("user").(models.UserInRedisOpenID)
}

// All ADs ...
func (ctrl ADUserController) AllAd(c *gin.Context) {
	if userID := getUserID(c); userID == "" {
		// c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty domain name"})
		// return
	}
	// var userRoles []string
	// if userID != "" {
	// }
	users, err := aduserModel.AllAd()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// All ...
func (ctrl ADUserController) All(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty domain name"})
		return
	}
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

// All ...
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
func (ctrl ADUserController) Find(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		adUser, adErr := aduserModel.GetOneUser(userID)
		if adErr != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
			return
		}
		//	jadUser, _ := json.Marshal(adUser)

		c.JSON(http.StatusOK, gin.H{"message": "User finded", "user": adUser})
	}
}

func (ctrl ADUserController) GetUserAdusers(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	role, err := aduserModel.GetAdusersRights(user)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role finded", "data": role})
}

func (ctrl ADUserController) GetUserByName(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	adUser, adErr := aduserModel.GetOneUserPropertys(user)
	if adErr != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User finded", "data": adUser})
}
