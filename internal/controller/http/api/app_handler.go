package api

import (
	"net/http"
	"strings"

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

func (h *Handler) getPerms(c *gin.Context) (entity.Permissions, bool) {
	permsVal, exists := c.Get("perms")
	if !exists {
		return entity.Permissions{}, false
	}

	perms, ok := permsVal.(entity.Permissions)
	return perms, ok
}

func (h *Handler) CurrentUser(c *gin.Context) {
	if !h.uc.IsAppInitialized() {
		c.JSON(http.StatusAccepted, gin.H{"message": "Приложение не инициализированно"})
		return
	}
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}

	adUser, adErr := h.uc.GetCurrentUser(perms)
	if adErr != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка получения данных пользователя", "error": adErr.Error()})
		return
	}
	//	jadUser, _ := json.Marshal(adUser)

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь найден", "user": adUser})

}

func (h *Handler) CurrentUserResources(c *gin.Context) {

	if userID := getUserID(c); userID != "" {
		resources, err := h.uc.GetCurrentUserResources(userID)

		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": resources})
		return
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid credentials"})
		return
	}
}

func (h *Handler) DomainList(c *gin.Context) {

	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(403, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	// userID := ""
	// if userID = getUserID(c); userID == "" {
	// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
	// 	return
	// }
	domainList := h.uc.GetDomainList(perms)
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
func (h *Handler) Ip(c *gin.Context) {

	ip := ReadUserIP(c.Request)

	c.HTML(http.StatusOK, "ip.html", gin.H{
		"Ip": ip,
	})

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
	roles, err := h.uc.GetAppRoles()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список ролей приложения", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *Handler) AppSections(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	roles, err := h.uc.GetAppSections()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список ролей приложения", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
func (h *Handler) AppDomains(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	domains := h.uc.DomainList()

	c.JSON(http.StatusOK, gin.H{"data": domains})
}
func (h *Handler) ComputerRMS(c *gin.Context) {
	isApi := h.IsApiRequest(c)
	if !isApi {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Этот запрос обслуживает только другие API"})
		return
	}
	domain := c.Param("domain")
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить домен.", "error": "Empty domain name"})
		return
	}
	computers, err := h.uc.GetComputerRMS(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Невозможно получить список коспьютеров с RMS", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, computers)
}

func (h *Handler) GetLocalAdmins(c *gin.Context) {

	computer := c.Param("computer")
	if computer == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Невозможно определить имя компьютера.", "error": "Empty computer name"})
		return
	}

	adminsStr := c.PostForm("administrators")
	if adminsStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":              "no administrators",
			"local_admins":        []string{},
			"local_admins_domain": []string{},
		})
		return
	}

	adminNames := strings.Split(adminsStr, ",")

	if len(adminNames) == 0 || (len(adminNames) == 1 && adminNames[0] == "") {
		c.JSON(http.StatusOK, gin.H{
			"status":              "empty",
			"local_admins":        []string{},
			"local_admins_domain": []string{},
		})
		return
	}

	var localAdmins []string       // Только имена после "/"
	var localAdminsDomain []string // Всё остальное

	prefix := "WinNT://BRNV/"
	for _, admin := range adminNames {
		cleanAdmin := strings.TrimSpace(admin)
		if after, ok := strings.CutPrefix(cleanAdmin, prefix); ok {
			// Есть префикс — ищем "/"
			if idx := strings.IndexByte(after, '/'); idx != -1 {
				// localAdmins: всё после "/"
				localAdmins = append(localAdmins, after[idx+1:])
			} else {
				// localAdminsDomain: всё после префикса (без "/")
				localAdminsDomain = append(localAdminsDomain, after)
			}
		} else {
			// Нет префикса — в localAdminsDomain
			localAdminsDomain = append(localAdminsDomain, cleanAdmin)
		}
	}

	if err := h.uc.ComputerLocalAdminsAudit(computer, localAdminsDomain, true); err != nil {
		h.log.Info(err)
	}
	if err := h.uc.ComputerLocalAdminsAudit(computer, localAdmins, false); err != nil {
		h.log.Info(err)
	}

	c.JSON(200, gin.H{
		"status":              "ok",
		"local_admins":        localAdmins,
		"local_admins_domain": localAdminsDomain,
	})
}
