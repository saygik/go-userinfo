package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OrgCodes(c *gin.Context) {
	user, err := h.uc.GetOrgCodes()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "ошибка предоставления списка кодов организаций", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}
