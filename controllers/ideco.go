package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type IDECOController struct{}

var idecoModel = new(models.IDECOModel)

// All ...
func (ctrl IDECOController) GetWiteList(c *gin.Context) {

	wlists, err := idecoModel.GetWiteList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get zals list from sharepoint server", "error": err.Error()})
		return
	}
	for _, wlist := range wlists {
		if wlist.Name == "Белый Список" {

			c.JSON(http.StatusOK, gin.H{"data": wlist.Urls})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": "{}"})
}
