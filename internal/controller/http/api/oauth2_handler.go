package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) GetOauth2Login(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
// }

func (h *Handler) RefreshOauth2(c *gin.Context) {

	var tokenForm struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid refresh token"})
		c.Abort()
		return
	}
	if tokenForm.RefreshToken == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "empty refresh token"})
		c.Abort()
		return
	}

	resp, err := h.hydra.IntrospectOAuth2Token(tokenForm.RefreshToken)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid refresh token", "err": err.Error()})
		c.Abort()
		return
	}
	if !resp.Active {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "refresh token is not active"})
		c.Abort()
		return
	}
	newToken, err := h.oAuth2.Refresh(tokenForm.RefreshToken)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid refresh token", "err": err.Error()})
		c.Abort()
		return
	}
	tt, err := h.hydra.IntrospectOAuth2Token(newToken.RefreshToken)
	_ = tt
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid new refresh token", "err": err.Error()})
		c.Abort()
		return
	}

	newToken.Expiry = time.Unix(*tt.Exp, 0)
	c.JSON(http.StatusOK, newToken)

}

func (h *Handler) LoginOauth2(c *gin.Context) {

	state := c.Query("state")

	loginURL := h.oAuth2.AuthCodeURL(state)

	c.JSON(http.StatusOK, gin.H{"data": loginURL})
}

func (h *Handler) LogoutOauth2(c *gin.Context) {
	var ticketForm struct {
		IdToken string `json:"id_token,omitempty"`
		Url     string `json:"url,omitempty"`
	}
	c.ShouldBindJSON(&ticketForm)
	url := h.hydra.LogoutURL() + "?post_logout_redirect_uri=" + ticketForm.Url + "&id_token_hint=" + ticketForm.IdToken

	c.JSON(http.StatusOK, gin.H{"data": url})
}

func (h *Handler) ExchangeToken(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusOK, gin.H{"error": "authorization code is empty"})
		return
	}

	accessToken, userInfo, err := h.oAuth2.Exchange(code)
	_ = userInfo
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	tt, err := h.hydra.IntrospectOAuth2Token(accessToken.RefreshToken)
	_ = tt
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	accessToken.Expiry = time.Unix(*tt.Exp, 0)

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
		"user":  userInfo,
	})
}

func (h *Handler) TokenValid(c *gin.Context) {

	token := h.ExtractToken(c.Request)

	if len(token) < 1 {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	resp, err := h.hydra.IntrospectOAuth2Token(token)
	if err != nil { //if any goes wrong
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}
	_ = resp
	if !resp.Active {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token is not active"})
		return
	}
	//To be called from GetUserID()
	c.Set("user", *resp.Sub)
}
