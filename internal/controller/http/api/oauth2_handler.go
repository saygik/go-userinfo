package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
}

func (h *Handler) PostLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
}

func (h *Handler) GetConsent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
}

func (h *Handler) PostConsent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
}
