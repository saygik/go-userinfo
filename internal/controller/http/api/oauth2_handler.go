package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) GetOauth2Login(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
// }

func (h *Handler) RefreshOauth2Authentik(c *gin.Context) {

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

func (h *Handler) LoginOauth2Authentik(c *gin.Context) {

	state := c.Query("state")

	loginURL := h.oAuth2Authentik.AuthCodeURL(state)

	c.JSON(http.StatusOK, gin.H{"data": loginURL})
}

func (h *Handler) LogoutOauth2Authentik(c *gin.Context) {
	var ticketForm struct {
		Url string `json:"url,omitempty"`
	}
	c.ShouldBindJSON(&ticketForm)
	url := h.oAuth2Authentik.LogOutURL()
	c.JSON(http.StatusOK, gin.H{"data": url})
}

func (h *Handler) ExchangeTokenAuthentik(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusOK, gin.H{"error": "authorization code is empty"})
		return
	}

	accessToken, userInfo, err := h.oAuth2Authentik.Exchange(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
		"user":  userInfo,
	})

}
