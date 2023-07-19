package models

import (
	"github.com/saygik/go-userinfo/db"
)

// UserModel ...
type ManualsModel struct{}

type OrgWithCodes struct {
	Id    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Tname string `db:"tname" json:"tname"`
	Key   string `db:"key" json:"code"`
}

// GLPI User find by Mail ...
func (m ManualsModel) AllOrgCodes() (orgs []OrgWithCodes, err error) {
	_, err = db.GetDB().Select(&orgs, "getOrgCodes")
	return orgs, err
}
