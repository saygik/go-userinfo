package forms

type UserActivityForm struct {
	User      string `form:"user" json:"user" binding:"required,email"`
	Ip        string `form:"ip" json:"ip"`
	Computer  string `form:"computer" json:"computer,omitempty"`
	Activiy   string `form:"activity" json:"activity,omitempty"`
	ActiviyIp string `form:"activityip" json:"activityip,omitempty"`
	Date      string `form:"date" json:"date,omitempty"`
}

type ScheduleTask struct {
	Id             string `db:"id" json:"id"`
	Idc            string `db:"idc" json:"idc" `
	Tip            string `db:"tip" json:"tip"`
	Status         string `db:"status" json:"status"`
	Title          string `db:"title" json:"title"`
	Start          string `db:"date_start" json:"start"`
	End            string `db:"date_end" json:"end"`
	Upn            string `db:"upn" json:"upn"`
	AllDay         string `db:"all_day" json:"all_day"`
	SendMattermost string `db:"send_mattermost" json:"send_mattermost"`
	Comment        string `db:"comment" json:"comment"`
}
type Schedule struct {
	Id        string `db:"id" json:"id"`
	Tip       string `db:"tip" json:"tip"`
	Name      string `db:"name" json:"name"`
	Domain    string `db:"domain" json:"domain"`
	Usergroup string `db:"usergroup" json:"usergroup"`
}

// SoftwareForm ...
type SoftwareForm struct {
	Id   int64  `form:"id" json:"id" binding:"required,number"`
	User string `form:"user" json:"user"`
}
