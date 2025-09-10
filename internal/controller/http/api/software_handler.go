package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) GetSoftwareUser(c *gin.Context) {
	userName := c.Param("username")

	user, err := h.uc.GetUserSoftwares(userName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления списка систем пользователя", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) GetSoftwares(c *gin.Context) {
	user, err := h.uc.GetSoftwares()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления списка систем", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})

}
func (h *Handler) GetSoftwaresUsers(c *gin.Context) {
	user, err := h.uc.GetSoftwaresUsers()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления списка систем", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})

}
func (h *Handler) AddSoftwareUser(c *gin.Context) {
	userName := c.Param("username")

	var softwareForm entity.SoftwareForm
	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&softwareForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить добавляемую систему: " + err.Error()})
		return
	}
	softwareForm.User = userName
	err = h.uc.AddUserSoftware(softwareForm)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка добавления системы пользователя", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "система добавлена"})
}

func (h *Handler) DelSoftwareUser(c *gin.Context) {
	id := c.Param("id")

	err := h.uc.DelUserSoftware(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка удаления системы пользователя", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "система удалена"})
}

func (h *Handler) GetSoftware(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить id сиситемы в запросе"})
		return
	}
	software, err := h.uc.GetSoftware(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "невозможно получить систему", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": software})
}

func (h *Handler) GetSoftwareUsers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить id сиситемы в запросе"})
		return
	}
	softUsers, err := h.uc.GetSoftwareUsers(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "невозможно получить список пользователей системы", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": softUsers})
}

func (h *Handler) AddUserToSoftware(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить id сиситемы в запросе"})
		return
	}
	var softwareForm entity.SoftUser
	err := c.ShouldBindJSON(&softwareForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить пользователя системы: " + err.Error()})
		return
	}

	userProperties, err := h.uc.AddOneSoftwareUser(id, softwareForm, userID)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка добавления пользователя к системе", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userProperties})
}

func (h *Handler) UpdateUserInSoftware(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить id сиситемы в запросе"})
		return
	}
	var userForm entity.SoftUser
	err := c.ShouldBindJSON(&userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Невозможно определить пользователя системы: " + err.Error()})
		return
	}

	userProperties, err := h.uc.UpdateOneSoftwareUser(userForm, userID)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка добавления пользователя к системе", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userProperties})
}
