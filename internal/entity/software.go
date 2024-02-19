package entity

type Software struct {
	Id             int64                    `db:"id" json:"id"`
	Name           string                   `db:"name" json:"name"`
	Ename          string                   `db:"ename" json:"company"`
	Comment        string                   `db:"comment" json:"comment"`
	Locations      string                   `db:"locations" json:"locations,omitempty"`
	Manufacture    string                   `db:"manufacture" json:"manufacture"`
	Description1   string                   `db:"description1" json:"description1"`
	Description2   string                   `db:"description2" json:"description2"`
	Murl           string                   `db:"murl" json:"manual_url"`
	Durl           string                   `db:"durl" json:"icon_url"`
	IsRecursive    int64                    `db:"is_recursive" json:"is_recursive"`
	Groups_id_tech int64                    `db:"groups_id_tech" json:"groups_id_tech"`
	Users_id_tech  int64                    `db:"users_id_tech" json:"users_id_tech"`
	Extauth        int64                    `db:"extauth" json:"ext_auth"`
	Clients        int64                    `db:"clients" json:"clients"`
	GroupName      string                   `db:"group_name" json:"group_name"`
	Admins         []map[string]interface{} `json:"tech_users"`
}

type SoftwareAdmins struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
