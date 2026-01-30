package entity

type UserIPComputer struct {
	Login    string `db:"login" json:"login"`
	Ip       string `db:"ip" json:"ip"`
	Computer string `db:"computer" json:"computer"`
	Rms      string `db:"rms" json:"rms"`
	IpDate   string `db:"ip_date" json:"ip_date"`
}

type ComputerUser struct {
	UPN         string `json:"upn,omitempty"`
	Company     string `json:"company,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Department  string `json:"department,omitempty"`
	Title       string `json:"title,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Telephone   string `json:"telephone,omitempty"`
	Computer    string `json:"computer,omitempty"`
	IP          string `json:"ip,omitempty"`
	LastDate    string `json:"last_date,omitempty"`
}
type ComputerProperties struct {
	OperatingSystem string `json:"operatingSystem,omitempty"`
	Description     string `json:"description,omitempty"`
}

// DomainComputer описывает компьютер домена, получаемый из MSSQL.
// Процедура GetComputerByDomain возвращает поля:
//   - computer   string
//   - ip         string
//   - last_date  string
//   - domain     string
//   - days_on    int
type DomainComputer struct {
	ID              int            `db:"id" json:"id"`
	Computer        string         `db:"computer" json:"computer"`
	IP              string         `db:"ip" json:"ip"`
	LastDate        string         `db:"last_date" json:"last_date"`
	Domain          string         `db:"domain" json:"domain"`
	DaysOn          int            `db:"days_on" json:"days_on"`
	Users           []ComputerUser `json:"users"`
	OperatingSystem string         `json:"operatingSystem,omitempty"`
	Description     string         `json:"description,omitempty"`
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
	Rms       string `form:"rms" json:"rms,omitempty"`
	Activiy   string `form:"activity" json:"activity,omitempty"`
	ActiviyIp string `form:"activityip" json:"activityip,omitempty"`
	Date      string `form:"date" json:"date,omitempty"`
}

type SimpleUser struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Company    string `json:"company,omitempty"`
	Department string `json:"department,omitempty"`
	Title      string `json:"title,omitempty"`
	Mail       string `json:"mail,omitempty"`
	Telephone  string `json:"telephone,omitempty"`
}
