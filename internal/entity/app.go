package entity

type IdName struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type AppResource struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Edit string `db:"edit" json:"edit"`
}
