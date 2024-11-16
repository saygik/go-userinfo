package entity

type HRPUser struct {
	Id      int    `json:"id"`
	FIO     string `json:"fio"`
	Dolg    string `json:"dolg"`
	Company string `json:"company"`
	Mero    string `json:"mero"`
	Date    string `json:"date"`
}
