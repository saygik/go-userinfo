package entity

type Schedule struct {
	Id             int                    `db:"id" json:"id"`
	Tip            int                    `db:"tip" json:"tip"`
	Name           string                 `db:"name" json:"name"`
	Private        bool                   `db:"private" json:"private"`
	Caption        string                 `db:"caption" json:"caption"`
	Edit           bool                   `db:"edit" json:"edit"`
	Usergroup      string                 `db:"usergroup" json:"usergroup"`
	Mattermost     string                 `db:"mattermost" json:"mattermost"`
	Available      bool                   `json:"available"`
	ScheduleUsers  []ScheduleUser         `json:"users"`
	ScheduleAdmins []IdName               `json:"admins"`
	ScheduleTasks  []ScheduleTaskCalendar `json:"tasks"`
}

type ExtendedProps struct {
	Id                 string `json:"id"`
	Title              string `json:"title"`
	SendMattermost     bool   `json:"sendMattermost"`
	NotificationSended bool   `json:"notificationSended"`
	Comment            string `json:"comment"`
	Status             int    `json:"status"`
	Tip                int    `json:"tip"`
	Company            string `json:"company"`
	Department         string `json:"department"`
	Mail               string `json:"mail"`
	TelephoneNumber    string `json:"telephoneNumber"`
	Mobile             string `json:"mobile"`
	Notfound           bool   `json:"notfound"`
}
type ScheduleTaskCalendar struct {
	Id            int           `json:"id"`
	Title         string        `json:"title"`
	Start         string        `json:"start"`
	End           string        `json:"end"`
	AllDay        bool          `json:"allDay"`
	ExtendedProps ExtendedProps `json:"extendedProps"`
}
type ScheduleTask struct {
	Id                 int    `db:"id" json:"id"`
	Idc                int    `db:"idc" json:"idc" `
	Tip                int    `db:"tip" json:"tip"`
	Status             int    `db:"status" json:"status"`
	Title              string `db:"title" json:"title"`
	Start              string `db:"start" json:"start"`
	End                string `db:"end" json:"end"`
	Upn                string `db:"upn" json:"upn"`
	AllDay             bool   `db:"allDay" json:"allDay"`
	SendMattermost     bool   `db:"sendMattermost" json:"sendMattermost"`
	NotificationSended bool   `db:"notificationSended" json:"notificationSended"`
	Comment            string `db:"comment" json:"comment"`
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
