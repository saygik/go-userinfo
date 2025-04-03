package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) AddOneSoftwareUser(form entity.SoftUser, editor string, currentTime string) (formRes entity.SoftUser, err error) {
	err = r.db.QueryRow("AddOneSoftwareUser_2 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10", form.Name, form.Id, form.Login, form.Comment, form.Fio, form.External, form.Mail, form.EndDate, editor, currentTime).Scan(&formRes.Id,
		&formRes.Name, &formRes.Login, &formRes.Fio, &formRes.External, &formRes.Comment, &formRes.EndDate, &formRes.Mail, &formRes.Sended, &formRes.Editor, &formRes.LastChanges)
	return formRes, err
}

func (r *Repository) UpdateOneSoftwareUser(form entity.SoftUser, editor string, currentTime string) error {
	res, err := r.db.Exec("UpdateOneSoftwareUser $1,$2,$3,$4,$5,$6,$7,$8,$9", form.Id, form.Login, form.Fio, form.Comment, form.EndDate, form.Mail, form.Sended, editor, currentTime)
	return sqlRowAffectedErrorWrapper(res, err)
}
