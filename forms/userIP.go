package forms

type UserActivityForm struct {
	User      string `form:"user" json:"user" binding:"required,email"`
	Ip        string `form:"ip" json:"ip"`
	Activiy   string `form:"activity" json:"activity,omitempty"`
	ActiviyIp string `form:"activityip" json:"activityip,omitempty"`
	Date      string `form:"date" json:"date,omitempty"`
}

type ScheduleTask struct {
	Id    string `db:"id" json:"id"`
	Idc   string `db:"idc" json:"idc" `
	Start string `db:"date_start" json:"start"`
	End   string `db:"date_end" json:"end"`
	Title string `db:"title" json:"title"`
	Upn   string `db:"upn" json:"upn"`
}
type Schedule struct {
	Id        string `db:"id" json:"id"`
	Tip       string `db:"tip" json:"tip"`
	Name      string `db:"name" json:"name"`
	Domain    string `db:"domain" json:"domain"`
	Usergroup string `db:"usergroup" json:"usergroup"`
}
