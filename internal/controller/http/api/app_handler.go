package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) CurrentUser(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		adUser, adErr := h.uc.GetUser(userID)
		if adErr != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": adErr.Error()})
			return
		}
		//	jadUser, _ := json.Marshal(adUser)

		c.JSON(http.StatusOK, gin.H{"message": "User finded", "user": adUser})
	}

}

func (h *Handler) CurrentUserResources(c *gin.Context) {
	if userID := getUserID(c); userID != "" {
		resources, err := h.uc.GetCurrentUserResources(userID)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": resources})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid credentials"})
		return
	}
}

func (h *Handler) DomainList(c *gin.Context) {
	userID := ""
	if userID = getUserID(c); userID == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get domains", "error": "Empty domain name"})
		return
	}
	domainList := h.uc.GetDomainList(userID)
	c.JSON(http.StatusOK, gin.H{"data": domainList})
}
