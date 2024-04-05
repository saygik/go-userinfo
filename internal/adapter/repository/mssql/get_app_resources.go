package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetAppResources() (roles []entity.IdName, err error) {
	_, err = r.db.Select(&roles, "GetAppResources")
	return roles, err
}
