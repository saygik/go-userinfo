package entity

type UserIPComputer struct {
	Login    string `db:"login" json:"login"`
	Ip       string `db:"ip" json:"ip"`
	Computer string `db:"computer" json:"computer"`
}

type IdNameAvatar struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}
