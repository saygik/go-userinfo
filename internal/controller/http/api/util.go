package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	if strings.Index(IPAddress, "::1") > -1 {
		return "127.0.0.1"
	}
	i := strings.Index(IPAddress, ":")
	if i > -1 {
		return IPAddress[:i]
	} else {
		return IPAddress
	}

}
func (h *Handler) clientIdFromRequestUrl(RequestUrl string) string {
	clientId := ""
	if strings.Index(RequestUrl, "client_id") > 1 {
		clientId = RequestUrl[strings.Index(RequestUrl, "client_id")+10:]
		if strings.Index(clientId, "&") > 1 {
			clientId = clientId[:strings.Index(clientId, "&")]
		}
	}
	return clientId
}

// ExtractToken ...
func (h *Handler) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	if len(bearToken) < 1 {
		return ""
	}
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (h *Handler) TokenValid(c *gin.Context) {

	token := h.ExtractToken(c.Request)

	if len(token) < 1 {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}
	user := ""
	if strings.HasPrefix(token, "ory_") {
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
		user = *resp.Sub
	} else {

		userInfo, err := h.oAuth2Authentik.IntrospectOAuth2Token(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token is not active"})
			return
		}
		user = userInfo.Sub
	}
	c.Set("user", user)

}

// UserFromToken ...
func (h *Handler) UserFromToken(c *gin.Context) {

	token := h.ExtractToken(c.Request)

	if len(token) < 1 {
		//Token either expired or not valid
		return
	}
	user := ""
	if strings.HasPrefix(token, "ory_") {
		resp, err := h.hydra.IntrospectOAuth2Token(token)
		if err != nil { //if any goes wrong
			return
		}
		_ = resp
		if !resp.Active {
			return
		}
		user = *resp.Sub
	} else {
		userInfo, err := h.oAuth2Authentik.IntrospectOAuth2Token(token)
		if err != nil {
			return
		}
		user = userInfo.Sub
	}
	//To be called from GetUserID()
	c.Set("user", user)
}
