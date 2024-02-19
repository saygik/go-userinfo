package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
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

	token, err := h.jwt.Login(loginForm.Email)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "User find in AD, but it not registred in this system", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен", "user": loginForm.Email, "token": token})
}

func (h *Handler) Logout(c *gin.Context) {

	au, err := h.jwt.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	err = h.jwt.DeleteAuth(au.AccessUUID)
	if err != nil { //if any goes wrong
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *Handler) TokenValid(c *gin.Context) {

	tokenAuth, err := h.jwt.ExtractTokenMetadata(c.Request)
	if err != nil {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	user, err := h.jwt.FetchAuth(tokenAuth)
	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	//To be called from GetUserID()
	c.Set("user", user)
}

func (h *Handler) Refresh(c *gin.Context) {
	var tokenForm entity.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form", "form": tokenForm})
		c.Abort()
		return
	}

	tokens, err := h.jwt.RefreshToken(tokenForm.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, tokens)
	}
}
