package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) AddOneSoftwareUser(form entity.SoftUser) (err error) {
	res, err := r.db.Exec("AddOneSoftwareUser $1,$2,$3,$4,$5,$6", form.Name, form.Id, form.Login, form.Comment, form.Fio, form.External)
	return sqlRowAffectedErrorWrapper(res, err)

}
