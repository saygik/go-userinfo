package entity

type UserIPComputer struct {
	Login    string `db:"login" json:"login"`
	Ip       string `db:"ip" json:"ip"`
	Computer string `db:"computer" json:"computer"`
}

type IdNameAvatar struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}

type UserActivity struct {
	Ip       string `db:"ip" json:"ip"`
	Activity string `db:"activity" json:"activity"`
	Date     string `db:"date" json:"date"`
}

type UserActivityForm struct {
	User      string `form:"user" json:"user" binding:"required,email"`
	Ip        string `form:"ip" json:"ip"`
	Computer  string `form:"computer" json:"computer,omitempty"`
	Activiy   string `form:"activity" json:"activity,omitempty"`
	ActiviyIp string `form:"activityip" json:"activityip,omitempty"`
	Date      string `form:"date" json:"date,omitempty"`
}
