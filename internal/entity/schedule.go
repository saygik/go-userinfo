package entity

type Schedule struct {
	Id             int            `db:"id" json:"id"`
	Tip            int            `db:"tip" json:"tip"`
	Name           string         `db:"name" json:"name"`
	Private        bool           `db:"private" json:"private"`
	Caption        string         `db:"caption" json:"caption"`
	Edit           bool           `db:"edit" json:"edit"`
	Usergroup      string         `db:"usergroup" json:"usergroup"`
	ScheduleUsers  []ScheduleUser `json:"users"`
	ScheduleAdmins []IdName       `json:"admins"`
	ScheduleTasks  []ScheduleTask `json:"tasks"`
}

type ScheduleTask struct {
	Id             int    `db:"id" json:"id"`
	Idc            int    `db:"idc" json:"idc" `
	Tip            int    `db:"tip" json:"tip"`
	Status         int    `db:"status" json:"status"`
	Title          string `db:"title" json:"title"`
	Start          string `db:"date_start" json:"start"`
	End            string `db:"date_end" json:"end"`
	Upn            string `db:"upn" json:"upn"`
	AllDay         bool   `db:"all_day" json:"all_day"`
	SendMattermost bool   `db:"send_mattermost" json:"send_mattermost"`
	Comment        string `db:"comment" json:"comment"`
}

type ScheduleUser struct {
	Name            string `json:"name"`
	AD              string `json:"ad"`
	DisplayName     string `json:"displayname,omitempty"`
	Company         string `json:"company,omitempty"`
	Title           string `json:"title,omitempty"`
	Department      string `json:"department,omitempty"`
	Mail            string `json:"mail,omitempty"`
	TelephoneNumber string `json:"telephonenumber,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
}
