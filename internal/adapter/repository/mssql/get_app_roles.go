package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetAppRoles() (roles []entity.IdName, err error) {
	_, err = r.db.Select(&roles, "GetAppRoles")
	return roles, err
}
