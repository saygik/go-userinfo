package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GlpiCurrentUser(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		user, err := h.uc.GetGlpiUser(userID)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}

}

func (h *Handler) GetGLPIUser(c *gin.Context) {
	userName := c.Param("username")
	userTechName := getUserID(c)
	user, err := h.uc.GetGlpiUserForTechnical(userName, userTechName)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Ошибка предоставления профиля пользователя GLPI", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}
