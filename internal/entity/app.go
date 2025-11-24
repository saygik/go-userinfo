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
