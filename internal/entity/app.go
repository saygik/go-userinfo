package entity

type Id struct {
	Id int `db:"id" json:"id"`
}

type IdName struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type IdNameType struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type int    `db:"type" json:"type"`
}

type AppResource struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Edit string `db:"edit" json:"edit"`
}

type IdNameFio struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Fio  string `db:"fio" json:"fio"`
}

type IdNameDescription struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type DomainAccess struct {
	Domain      string `db:"domain" json:"domain"`
	AccessLevel string `db:"access_level" json:"access_level"`
}
type Permissions struct {
	User         string          `json:"user"`
	Roles        []string        `json:"roles"`
	IsSysAdmin   bool            `json:"is_sysadmin"`
	IsAdmin      bool            `json:"is_admin"`
	IsTech       bool            `json:"is_tech"`
	HomeDomain   string          `json:"home_domain"`
	AdminDomains map[string]bool `json:"admin_domains"`
	TechDomains  map[string]bool `json:"tech_domains"`
	UserDomains  map[string]bool `json:"user_domains"`
	AllDomains   bool            `json:"all_domains"`
	Sections     []string        `json:"sections"`
}
