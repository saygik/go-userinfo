package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Wlist(c *gin.Context) {

	list, err := h.uc.IUTMGetWlist()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)

}
