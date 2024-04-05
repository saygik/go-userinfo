package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetAppGroups() (roles []entity.IdName, err error) {
	_, err = r.db.Select(&roles, "GetAppGroups")
	return roles, err
}
