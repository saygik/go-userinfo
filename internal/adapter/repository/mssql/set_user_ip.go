package mssql

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SetUserIp(form entity.UserActivityForm) (msgResponce string, err error) {
	//	_, err = db.GetDB().Exec("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy)
	err = r.db.QueryRow("SetUserIPActivityComputerRms $1,$2,$3,$4,$5", form.User, form.Ip, form.Computer, form.Activiy, form.Rms).Scan(&msgResponce)

	return msgResponce, err
}
