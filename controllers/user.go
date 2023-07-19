package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/forms"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type UserController struct{}

var userModel = new(models.UserModel)

// Login ...
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm
	err := c.ShouldBindJSON(&loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Неправильные данные"})
		return
	}
	var adModel = new(models.ADUserModel)
	_, _, adErr := adModel.Authenticate(loginForm)
	if adErr != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Неправильный логин или пароль", "error": adErr.Error()})
		return
	}

	token, err := userModel.Login(loginForm.Email)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "User find in AD, but it not registred in this system", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен", "user": loginForm.Email, "token": token})
}

// Logout ...
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

func (ctrl UserController) LoginOauth(c *gin.Context) {
	// cd := c.Request.URL.RawQuery
	// fmt.Printf(cd)
	redirect_uri := c.Query("redirect_uri")
	state := c.Query("state")
	url := "https://adss.brnv.rw/sso/oauth/7dd4815e19af5fbea99a290b134b7e493569ea13/authorize?client_id=4OcmRaXcBIsoTehRDcF5fYO3N&response_type=code&scope=openid+profile+email&redirect_uri=" + redirect_uri + "&state=" + state
	//url := "https://adss.brnv.rw/sso/oauth/7dd4815e19af5fbea99a290b134b7e493569ea13/authorize?client_id=4OcmRaXcBIsoTehRDcF5fYO3N&response_type=code&redirect_uri=http://172.28.7.203:8080/callback&scope=openid+profile+email"

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctrl UserController) LoginUser(c *gin.Context) {
	user := getUser(c)
	userInfo, err := authModel.GetUserInfo(user.AccessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}
