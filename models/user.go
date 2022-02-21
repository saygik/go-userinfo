package models

import (
	"fmt"
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
func (m UserModel) All() (users []User, err error) {
	sql := fmt.Sprintf(
		`SELECT  login, ip from UserIP`)
	rows, err := db.GetDB().Query(sql)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		// In each step, scan one row
		var user User
		err = rows.Scan(&user.Login, &user.Ip)
		if err != nil {
			return users, err
		}
		// and append it to the array
		users = append(users, user)
	}
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
