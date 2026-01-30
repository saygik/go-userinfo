package api

import (
	"fmt"
	"net/http"
	"strings"
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

// LastComputers возвращает компьютеры домена с последними логинами пользователей
// и дополнительными AD‑свойствами из Redis.
func (h *Handler) LastComputers(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен.", "error": "Empty domain name"})
		return
	}

	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить пользователя.", "error": "Empty user id"})
		return
	}

	computers, err := h.uc.GetADLastComputers(domain, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список компьютеров домена", "error": err.Error()})
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

// func (h *Handler) ADUserAddGroup(c *gin.Context) {
// 	user := "say@brnv.rw"
// 	group := "CN=as-test-group,OU=ИВЦ2АС,OU=Пользователи,OU=ИВЦ,OU=_Служебные записи,DC=brnv,DC=rw"
// 	err := h.uc.ADUserAddGroup(user, group)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно добавить пользователя в группу.", "error": "Failed to add user to group"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": "OK"})
// }

//	func (h *Handler) ADUserDelGroup(c *gin.Context) {
//		user := "say@brnv.rw"
//		group := "CN=as-test-group,OU=ИВЦ2АС,OU=Пользователи,OU=ИВЦ,OU=_Служебные записи,DC=brnv,DC=rw"
//		err := h.uc.ADUserDelGroup(user, group)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно удалить пользователя из группы.", "error": "Failed to delete user from group"})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"data": "OK"})
//	}
func (h *Handler) SwitchUserGroupInternet(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")

	type inet struct {
		Group       string `json:"group"`       // Группа: whitelist, full, tech или ""
		IsTemporary bool   `json:"isTemporary"` // Временное изменение (по умолчанию false)
		Days        int    `json:"days"`        // Количество суток (1 = завтра в 8:00, 2 = послезавтра в 8:00 и т.д., используется только если isTemporary = true)
	}
	var idForm inet

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить новую группу интернета: " + err.Error()})
		return
	}

	group := idForm.Group
	isTemporary := idForm.IsTemporary
	days := idForm.Days

	// Валидация days
	if isTemporary {
		if days < 1 {
			days = 1 // По умолчанию 1 сутки (завтра в 8:00)
		}
	}

	err = h.uc.SwitchUserGroupInternet(userID, user, group, isTemporary, days)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно добавить пользователя в группу", "error": err.Error()})
		return
	}

	message := "Группа успешно изменена"
	if isTemporary {
		dayWord := "суток"
		if days == 1 {
			dayWord = "сутки"
		}
		message = fmt.Sprintf("Группа временно изменена на %d %s (возврат в 8:00 утра)", days, dayWord)
	}

	c.JSON(http.StatusOK, gin.H{"data": "OK", "message": message})
}

// GetTemporaryGroupChange получает информацию о временном изменении группы пользователя
func (h *Handler) GetTemporaryGroupChange(c *gin.Context) {
	user := c.Param("username")

	change, err := h.uc.GetTemporaryGroupChange(user)
	if err != nil {
		// Если запись не найдена, возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Временное изменение не найдено", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": change})
}

// getDomainFromUserName извлекает домен из имени пользователя (user@domain)
func getDomainFromUserName(user string) string {
	parts := strings.Split(user, "@")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// DeleteTemporaryGroupChange удаляет временное изменение группы и восстанавливает предыдущую группу
func (h *Handler) DeleteTemporaryGroupChange(c *gin.Context) {
	userID := getUserID(c)
	user := c.Param("username")

	// Проверяем права доступа (только администраторы домена могут удалять)
	domain := getDomainFromUserName(user)
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат имени пользователя"})
		return
	}

	// Проверяем права (только администраторы домена могут удалять)
	permissionGroupDomainAdmins := "Администраторы пользователей домена"
	if err := h.uc.UserInDomainGroup2(userID, permissionGroupDomainAdmins, domain); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "У Вас нет прав на эту операцию"})
		return
	}

	err := h.uc.DeleteTemporaryGroupChange(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении временного изменения", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "OK", "message": "Временное изменение удалено, группа восстановлена"})
}
