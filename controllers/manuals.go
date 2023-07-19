package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type ManualsController struct{}

var manualsModel = new(models.ManualsModel)

func (ctrl ManualsController) AllOrgCodes(c *gin.Context) {

	results, err := manualsModel.AllOrgCodes()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get orgs ", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}
