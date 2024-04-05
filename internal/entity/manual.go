package entity

type OrgWithCodes struct {
	Id    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Tname string `db:"tname" json:"tname"`
	Key   string `db:"key" json:"code"`
}
