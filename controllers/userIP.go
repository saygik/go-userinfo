package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/forms"
	"github.com/saygik/go-userinfo/models"
)

// UserController ...
type UserIPController struct{}

var userIPModel = new(models.UserIPModel)
var DefaultDomain string

func (ctrl UserIPController) All(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		domain = DefaultDomain
	}
	results, err := userIPModel.All(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}
func (ctrl UserIPController) GetSchedule(c *gin.Context) {
	schedule := c.Param("id")
	if schedule == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указан ID расписания"})
		return
	}
	results, err := userIPModel.Schedule(schedule)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get schedule", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}
func (ctrl UserIPController) GetScheduleTasks(c *gin.Context) {
	schedule := c.Param("idc")
	if schedule == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Не указан ID расписания"})
		return
	}
	results, err := userIPModel.AllScheduleTasks(schedule)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})

}

func (ctrl UserIPController) DelScheduleTask(c *gin.Context) {
	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {

		ra, err := userIPModel.DelScheduleTask(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Schedule event could not be deleted", "error": err.Error()})
			return
		}
		if ra == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Schedule event does not exist"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
	}
}
func (ctrl UserIPController) AddScheduleTask(c *gin.Context) {
	//decoder := json.NewDecoder(c.Request.Body)
	var scheduleForm forms.ScheduleTask
	//err := decoder.Decode(&scheduleForm)
	err := c.ShouldBindJSON(&scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get schedule info from request", "error": err.Error()})
		return

	}

	msgResponce, err := userIPModel.AddScheduleTask(scheduleForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not add schedule", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": msgResponce})

}

func (ctrl UserIPController) UpdateScheduleTask(c *gin.Context) {
	id := c.Param("id")
	if id, err := strconv.ParseInt(id, 10, 64); err == nil {
		var scheduleForm forms.ScheduleTask
		err := c.ShouldBindJSON(&scheduleForm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get schedule info from request", "error": err.Error()})
			return
		}
		scheduleForm.Id = strconv.FormatInt(id, 10)
		ra, err := userIPModel.UpdateScheduleTask(scheduleForm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Событие календаря не изменено: ", "error": err.Error()})
			return
		}
		if ra == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Событие календаря для изменения не существует"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Событие календаря изменено"})

	} else {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "error": err.Error()})
	}

}

func (ctrl UserIPController) SetIp(c *gin.Context) {
	var userForm forms.UserActivityForm
	err := c.ShouldBindQuery(&userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": err.Error()})
		return

	}
	domain := strings.Split(userForm.User, "@")[1]

	if !ad.Domains[domain] {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set users ip", "error": "Non served domain"})
		return
	}

	userForm.Ip = ReadUserIP(c.Request)
	//	activity := query.Get("activity")
	if userForm.Activiy == "" {
		userForm.Activiy = "login"
	}
	var msgResponce string
	msgResponce, err = userIPModel.SetUserIp(userForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not set user ip", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msgResponce})
}

func (ctrl UserIPController) GetUserByName(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	userWithActivity, err := userIPModel.GetUserByName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users ip`s", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userWithActivity})
	// ipAddr:=ReadUserIP(c.Request)
	//var hostNames []string
	// hostNames,_=ReadUserName(ipAddr)
	//c.JSON(http.StatusOK, gin.H{"data": ipAddr, "host_names": hostNames})

}

func (ctrl UserIPController) GetUserWeekActivity(c *gin.Context) {
	user := c.Param("username")
	if user == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Could not get user info from request", "error": ""})
		return
	}
	if !isEmailValid(user) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid username. The name must include the domain in the format: user@domain", "error": ""})
		return
	}
	data, err := userIPModel.GetUserWeekActivity(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get users activity from server", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})

}
