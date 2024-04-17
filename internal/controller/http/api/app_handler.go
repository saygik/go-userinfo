package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

// getUserID ...
func getUserID(c *gin.Context) (userID string) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	_, isExist := c.Get("user")
	if isExist {
		return c.MustGet("user").(string)
	}
	return ""
}

func (h *Handler) CurrentUser(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		adUser, adErr := h.uc.GetCurrentUser(userID, userID)
		if adErr != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
			return
		}
		//	jadUser, _ := json.Marshal(adUser)

		c.JSON(http.StatusOK, gin.H{"message": "User finded", "user": adUser})
	}

}

func (h *Handler) CurrentUserResources(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		resources, err := h.uc.GetCurrentUserResources(userID)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": resources})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid credentials"})
		return
	}
}

func (h *Handler) DomainList(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	domainList := h.uc.GetDomainList(userID)
	c.JSON(http.StatusOK, gin.H{"data": domainList})
}

func (h *Handler) SetIp(c *gin.Context) {
	var userForm entity.UserActivityForm
	err := c.ShouldBindQuery(&userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": err.Error()})
		return
	}
	userForm.Ip = ReadUserIP(c.Request)
	msgResponce, err := h.uc.SetUserIp(userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set user ip", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msgResponce})

}

func (h *Handler) AppResources(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	resources, err := h.uc.GetAppResources(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список ресурсов приложения", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resources})
}

func (h *Handler) AppRoles(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	roles, err := h.uc.GetAppRoles(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список ролей приложения", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *Handler) AppGroups(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	groups, err := h.uc.GetAppGroups(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список групп пользователей приложения", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": groups})
}
