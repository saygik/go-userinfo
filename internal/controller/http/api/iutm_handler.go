package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Wlist(c *gin.Context, listName string) {
	list, err := h.uc.IUTMGetList(listName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains ", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// func (h *Handler) Wlist2(c *gin.Context) {

// 	list, err := h.uc.IUTMGetList("wlist2")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains ", "error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, list)

// }
// func (h *Handler) Wlistivc(c *gin.Context) {

// 	list, err := h.uc.IUTMGetList("wlistivc")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains ", "error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, list)

// }

// func (h *Handler) Blist(c *gin.Context) {

// 	list, err := h.uc.IUTMGetList("blist")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains ", "error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, list)

// }

func (h *Handler) AllLists(c *gin.Context) {

	list, err := h.uc.IUTMGetAllLists()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get all domains ", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)

}
