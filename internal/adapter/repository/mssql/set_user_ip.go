package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SetUserIp(form entity.UserActivityForm) (msgResponce string, err error) {
	//	_, err = db.GetDB().Exec("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy)
	err = r.db.QueryRow("SetUserIPActivityComputer $1,$2,$3,$4", form.User, form.Ip, form.Computer, form.Activiy).Scan(&msgResponce)

	return msgResponce, err
}
