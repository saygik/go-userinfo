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

//Login ...
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm
	err := c.ShouldBindJSON(&loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form"})
		return
	}
	var adModel = new(models.ADUserModel)
	_, _, adErr := adModel.Authenticate(loginForm)
	if adErr != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
		return
	}

	token, err := userModel.Login(loginForm.Email)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "User find in AD, but it not registred in this system", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User signed in", "user": loginForm.Email, "token": token})
}

//Logout ...
func (ctrl UserController) Logout(c *gin.Context) {

	au, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
