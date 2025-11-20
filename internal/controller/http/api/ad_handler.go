package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) Users(c *gin.Context) {
	start := time.Now()
	defer func() {
		observeRequest(time.Since(start), c.Writer.Status())
	}()

	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty domain name"})
		return
	}
	users, err := h.uc.GetADUsers(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}

func (h *Handler) PUsers(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain users", "error": "Empty domain name"})
		return
	}
	users, err := h.uc.GetADUsersPublicInfo(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}

func (h *Handler) GetAdCounts(c *gin.Context) {
	users, computers, err := h.uc.GetAdCounts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains counts", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "computers": computers})
}

// All ADs Computers ...
func (h *Handler) Computers(c *gin.Context) {

	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domains computers", "error": "Empty domain name"})
		return
	}
	computers, err := h.uc.GetADComputers(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains computers", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": computers})
}

func (h *Handler) User(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domain user", "error": "Empty domain name"})
		return
	}

	adUser, adErr := h.uc.GetUserADPropertys(user, userID)
	if adErr != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User finded", "data": adUser})
}
func (h *Handler) UserSimple(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.JSON(http.StatusOK, entity.SimpleUser{Name: "--", Department: "--"})
		//		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}

	adUser, err := h.uc.GetUserADPropertysSimple(user)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "Invalid user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adUser)
}

func (h *Handler) GetUserADActivity(c *gin.Context) {
	userName := c.Param("username")
	userTechName := getUserID(c)

	activity, err := h.uc.GetUserADActivity(userName, userTechName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления активности пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": activity})
}

func (h *Handler) GetUserMailboxPermissions(c *gin.Context) {
	userName := c.Param("username")
	userTechName := getUserID(c)

	activity, err := h.uc.GetUserMailboxPermissions(userName, userTechName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления делегированных почтовых ящиков для пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": activity})
}

func (h *Handler) UpdateUserAvatar(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var avatarForm entity.AvatarForm
	err := c.ShouldBindJSON(&avatarForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно получить имя аватара из запроса", "error": err.Error()})
		return
	}
	err = h.uc.SetUserAvatar(userID, user, avatarForm.Avatar)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить аватар пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.SetUserRole(userID, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) AddUserGroup(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить группу системы: " + err.Error()})
		return
	}
	err = h.uc.AddUserGroup(userID, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить группу пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}
func (h *Handler) AddUserRole(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.AddUserRole(userID, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) DelUserGroup(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить группу системы: " + err.Error()})
		return
	}
	err = h.uc.DelUserGroup(userID, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить группу пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) DelUserRole(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.DelUserRole(userID, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

// All users in group...
func (h *Handler) GroupUsers(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен пользователя.", "error": "Empty domain name"})
		return
	}
	group := c.Param("group")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить группу пользователя.", "error": "Empty group name"})
		return
	}
	users, err := h.uc.GetADGroupUsers(domain, group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно определить пользователей группы.", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
