package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) GetOauth2Login(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
// }

func (h *Handler) RefreshOauth2Authentik(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid refresh token"})
		c.Abort()
		return
	}

	if refreshToken == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "empty refresh token"})
		c.Abort()
		return
	}
	newToken, err := h.oAuth2Authentik.ExchangeRefreshToAccessToken(refreshToken)

	//	resp, err := h.hydra.IntrospectOAuth2Token(tokenForm.RefreshToken)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid refresh token", "err": err.Error()})
		c.Abort()
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newToken.RefreshToken,
		Path:     "/",
		Domain:   "brnv.rw", // or leave "" for host-only same-origin
		Secure:   true,      // required with SameSite=None
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	newToken.RefreshToken = ""
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
	//c.SetCookie("refresh_token", accessToken.RefreshToken, 30*24*60*60, "/", "", false, true) // adjust TTL
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    accessToken.RefreshToken,
		Path:     "/",
		Domain:   "brnv.rw", // or leave "" for host-only same-origin
		Secure:   true,      // required with SameSite=None
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
		"user":  userInfo,
	})

}
