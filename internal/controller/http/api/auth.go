package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	var loginForm entity.LoginForm
	err := c.ShouldBindJSON(&loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Неправильные данные"})
		return
	}

	_, _, adErr := h.uc.Authenticate(loginForm)
	if adErr != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Неправильный логин или пароль", "error": adErr.Error()})
		return
	}

	// token, err := userModel.Login(loginForm.Email)
	// if err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{"message": "User find in AD, but it not registred in this system", "error": err.Error()})
	// 	return
	// }

	//c.JSON(http.StatusOK, gin.H{"message": "Вход произведен", "user": loginForm.Email, "token": token})
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен", "user": loginForm.Email})
}

func (h *Handler) Logout(c *gin.Context) {

}
