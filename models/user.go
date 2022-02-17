package models

import (
	"fmt"
	"github.com/saygik/go-userinfo/db"
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
