package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
