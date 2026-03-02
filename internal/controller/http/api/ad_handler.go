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

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	users, err := h.uc.GetADUsers(perms)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}

func (h *Handler) PUsers(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	users, err := h.uc.GetADUsersPublicInfo(perms)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}

func (h *Handler) GetAdCounts(c *gin.Context) {
	domain := c.Param("domain")

	var (
		users     int
		computers int
		err       error
	)

	if domain == "" || domain == "-" || domain == "все домены" {
		// Старое поведение: считаем по всем доменам
		users, computers, err = h.uc.GetAdCounts()
	} else {
		// Новое поведение: считаем только по одному домену
		users, computers, err = h.uc.GetAdCountsDomain(domain)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains counts", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "computers": computers})
}

// All ADs Computers ...
func (h *Handler) Computers(c *gin.Context) {

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	computers, err := h.uc.GetADComputers(perms)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains computers", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": computers})
}

// ComputersVersions возвращает количество компьютеров домена, сгруппированное по версии ОС (operatingSystemVersion).
// Включает также человекочитаемую версию Windows (operatingSystemVersionHuman, например 24H2).
func (h *Handler) ComputersVersions(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен.", "error": "Empty domain name"})
		return
	}

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	data, err := h.uc.GetADComputersVersions(domain, perms)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить версии компьютеров домена", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// ComputersOSFamily возвращает количество компьютеров домена, сгруппированное по семейству ОС (OperatingSystemFamily).
func (h *Handler) ComputersOSFamily(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен.", "error": "Empty domain name"})
		return
	}

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	data, err := h.uc.GetADComputersOSFamily(domain, perms)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить статистику по семействам ОС домена", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// LastComputers возвращает компьютеры домена с последними логинами пользователей
// и дополнительными AD‑свойствами из Redis.
func (h *Handler) LastComputers(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен.", "error": "Empty domain name"})
		return
	}

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	computers, err := h.uc.GetADLastComputers(domain, perms)
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

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	adUser, adErr := h.uc.GetUserADPropertys(user, perms)
	if adErr != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
		return
	}
	roles, err := h.uc.GetUserRoles(user)
	if err == nil {
		adUser["roles"] = roles
	}
	sections, err := h.uc.GetSections(user)
	if err == nil {
		adUser["sections"] = sections
	}
	domains, err := h.uc.GetDomainAccess(user)
	if err == nil {
		adUser["domains"] = domains
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
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	activity, err := h.uc.GetUserADActivity(userName, perms)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления активности пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": activity})
}

func (h *Handler) GetUserMailboxPermissions(c *gin.Context) {
	userName := c.Param("username")
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	activity, err := h.uc.GetUserMailboxPermissions(userName, perms)
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

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}

	err = h.uc.SetUserRole(perms, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) AddUserSection(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.IdNameDescription

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	section, err := h.uc.AddUserSection(perms, user, "SECTION:"+idForm.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, section)
}

func (h *Handler) AddUserDomainRole(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.DomainAccess

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	domain, operation, err := h.uc.AddUserDomainRole(perms, user, idForm.Domain, idForm.AccessLevel)
	_ = operation
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}

	switch operation {
	case "INSERTED":
		c.JSON(http.StatusCreated, domain)
	case "UPDATED":
		c.JSON(http.StatusOK, domain)
	default:
		c.JSON(http.StatusOK, domain)
	}

}

func (h *Handler) DelUserDomainRole(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.DomainAccess

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.DelUserDomainRole(perms, user, idForm.Domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) AddUserRole(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	role, err := h.uc.AddUserRole(perms, user, idForm.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (h *Handler) DelUserSection(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.IdNameDescription

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.DelUserSection(perms, user, "SECTION:"+idForm.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно обновить роль пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h *Handler) DelUserRole(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	user := c.Param("username")
	var idForm entity.IdName

	err := c.ShouldBindJSON(&idForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить роль системы: " + err.Error()})
		return
	}
	err = h.uc.DelUserRole(perms, user, idForm.Id)
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
