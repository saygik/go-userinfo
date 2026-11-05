package entity

type Schedule struct {
	Id             int                    `db:"id" json:"id"`
	Tip            int                    `db:"tip" json:"tip"`
	Name           string                 `db:"name" json:"name"`
	Private        bool                   `db:"private" json:"private"`
	Caption        string                 `db:"caption" json:"caption"`
	Edit           bool                   `db:"edit" json:"edit"`
	Mattermost     string                 `db:"mattermost" json:"mattermost"`
	Available      bool                   `json:"available"`
	Managers       ScheduleAccess         `json:"managers"`
	Viewers        ScheduleAccess         `json:"viewers"`
	ScheduleUsers  []ScheduleUser         `json:"users"`
	ScheduleAdmins []IdName               `json:"admins"`
	ScheduleTasks  []ScheduleTaskCalendar `json:"tasks"`
}

// ScheduleAccess описывает выданные календарю права. Пользователи хранятся
// по UPN, группы — по идентификатору и имени группы GLPI.
type ScheduleAccess struct {
	Users  []string              `json:"users"`
	Groups []ScheduleAccessGroup `json:"groups"`
}

type ScheduleAccessGroup struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ScheduleAccessUserRow и ScheduleAccessGroupRow — строки таблиц доступа,
// которые репозиторий отдаёт в usecase.
type ScheduleAccessUserRow struct {
	Name string `db:"name" json:"name"`
	Role int    `db:"role" json:"role"`
}

type ScheduleAccessGroupRow struct {
	GroupId int `db:"group_id" json:"group_id"`
	Role    int `db:"role" json:"role"`
}

// ScheduleForm — форма создания и редактирования календаря в админке.
type ScheduleForm struct {
	Name       string         `json:"name"`
	Tip        int            `json:"tip"`
	Mattermost string         `json:"mattermost"`
	Private    bool           `json:"private"`
	Caption    string         `json:"caption"`
	Managers   ScheduleAccess `json:"managers"`
	Viewers    ScheduleAccess `json:"viewers"`
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
