package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) AddOneSoftwareUser(form entity.SoftUser) (formRes entity.SoftUser, err error) {
	err = r.db.QueryRow("AddOneSoftwareUser_2 $1,$2,$3,$4,$5,$6,$7,$8", form.Name, form.Id, form.Login, form.Comment, form.Fio, form.External, form.Mail, form.EndDate).Scan(&formRes.Id,
		&formRes.Name, &formRes.Login, &formRes.Fio, &formRes.External, &formRes.Comment, &formRes.EndDate, &formRes.Mail, &formRes.Sended)
	return formRes, err
}
