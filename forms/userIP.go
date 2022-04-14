package forms

type UserActivityForm struct {
	User      string `form:"user" json:"user" binding:"required,email"`
	Ip        string `form:"ip" json:"ip"`
	Activiy   string `form:"activity" json:"activity,omitempty"`
	ActiviyIp string `form:"activityip" json:"activityip,omitempty"`
	Date      string `form:"date" json:"date,omitempty"`
}
