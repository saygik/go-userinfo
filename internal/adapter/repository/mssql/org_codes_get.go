package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetOrgCodes() (orgs []entity.OrgWithCodes, err error) {
	_, err = r.db.Select(&orgs, "getOrgCodes")
	return orgs, err
}
