package entity

import "strings"

type SoftwareGroup struct {
	Id            int64  `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	GroupMatt     string `db:"group_matt_channel" json:"group_matt_channel"`
	GroupCalendar int    `db:"group_calendar" json:"group_calendar"`
}

type Software struct {
	Id           int64                    `db:"id" json:"id"`
	Name         string                   `db:"name" json:"name"`
	Ename        string                   `db:"ename" json:"company"`
	Login        string                   `db:"login" json:"login,omitempty"`
	Comment      string                   `db:"comment" json:"comment"`
	Locations    string                   `db:"locations" json:"locations,omitempty"`
	Manufacture  string                   `db:"manufacture" json:"manufacture"`
	Description1 string                   `db:"description1" json:"description1"`
	Description2 string                   `db:"description2" json:"description2"`
	Murl         string                   `db:"murl" json:"manual_url"`
	Durl         string                   `db:"durl" json:"icon_url"`
	IsRecursive  int64                    `db:"is_recursive" json:"is_recursive"`
	Admins       []map[string]interface{} `json:"tech_users"`
	Groups       []SoftwareGroup          `json:"groups"`
}

//Clients      int64                    `db:"clients" json:"clients"`
//	Users_id_tech int64  `db:"users_id_tech" json:"users_id_tech"`
//	Extauth       int64  `db:"extauth" json:"ext_auth"`
// GroupMatt        string                   `db:"group_matt_channel" json:"group_matt_channel"`
// GroupCalendar    int                      `db:"group_calendar" json:"group_calendar"`
// Groups_id_tech_s string                   `db:"groups_id_tech_s" json:"groups_id_tech_s"` // "1,2,5"
// GroupNames_s     string                   `db:"group_names_s" json:"group_names_s"`       // "Group1, Group2"
// Groups_id_tech   []int64                  `db:"-" json:"groups_id_tech"`
// GroupNames       []string                 `db:"-" json:"group_names"`

func (s *Software) GroupNames() string {
	if len(s.Groups) == 0 {
		return ""
	}

	names := make([]string, 0, len(s.Groups))
	for _, group := range s.Groups {
		names = append(names, group.Name)
	}
	return strings.Join(names, ", ")
}
func (s *Software) HasTechGroup(groupID int64) bool {
	for _, group := range s.Groups {
		if group.Id == groupID {
			return true
		}
	}
	return false
}

type SoftwareAdmins struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type SoftwareForm struct {
	Id   int64  `form:"id" json:"id" binding:"required,number"`
	User string `form:"user" json:"user"`
}

type SoftUser struct {
	Id          int64                  `db:"id" json:"id"`
	IdSoft      int                    `db:"id_soft" json:"id_soft"`
	SoftName    string                 `db:"soft_name" json:"soft_name,omitempty"`
	Name        string                 `db:"user_name" json:"name"`
	Login       string                 `db:"user_login" json:"login,omitempty"`
	Comment     string                 `db:"user_comment" json:"comment,omitempty"`
	Fio         string                 `db:"fio" json:"fio,omitempty"`
	External    bool                   `db:"external" json:"external"`
	EndDate     string                 `db:"enddate" json:"enddate,omitempty"`
	Mail        string                 `db:"mail" json:"mail,omitempty"`
	Sended      bool                   `db:"sended" json:"sended"`
	Editor      string                 `db:"editor" json:"editor"`
	LastChanges string                 `db:"last_changes" json:"last_changes"`
	Propertys   map[string]interface{} `json:"props"`
}
