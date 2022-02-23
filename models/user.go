package models

import (
	"github.com/saygik/go-userinfo/db"
	"github.com/saygik/go-userinfo/forms"
)

//UserModel ...
type UserModel struct{}

type User struct {
	Login string `db:"login" json:"login"`
	Ip    string `db:"ip" json:"ip"`
}

//GLPI User find by Mail ...
func (m UserModel) All(domain string) (users []User, err error) {
	_, err = db.GetDB().Select(&users, "GetAllUserIPByDomain $1", domain)
	return users, err
}
func (m UserModel) SetUserIp(form forms.UserActivityForm) (msgResponce string, err error) {
	//	_, err = db.GetDB().Exec("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy)
	err = db.GetDB().QueryRow("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy).Scan(&msgResponce)

	return msgResponce, err
}
func (m UserModel) GetUserByName(user string) (form forms.UserActivityForm, err error) {
	err = db.GetDB().QueryRow("GetUserByName $1", user).Scan(&form.User, &form.Ip, &form.Activiy, &form.ActiviyIp, &form.Date)

	return form, err
}
