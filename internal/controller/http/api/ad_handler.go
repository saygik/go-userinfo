package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Users(c *gin.Context) {
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
