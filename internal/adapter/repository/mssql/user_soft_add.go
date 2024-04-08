package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) AddOneUserSoftware(form entity.SoftwareForm) (err error) {
	res, err := r.db.Exec("AddOneUserSoftware $1,$2", form.User, form.Id)
	return sqlRowAffectedErrorWrapper(res, err)
}
