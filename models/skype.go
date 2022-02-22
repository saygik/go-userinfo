package models

import (
	"github.com/saygik/go-userinfo/db"
)

//UserModel ...
type SkypeModel struct{}

type UserPresence struct {
	Userathost  string `db:"userathost" json:"userathost"`
	Presence    string `db:"presence" json:"presence"`
	Lastpubtime string `db:"lastpubtime" json:"lastpubtime"`
}

//GLPI User find by Mail ...
func (m SkypeModel) AllPresences() (userPresences []UserPresence, err error) {

	rows, err := db.GetDBSkype().Query("Skype_GetAvailability")
	if err != nil {
		return userPresences, err
	}
	for rows.Next() {
		// In each step, scan one row
		var userPresence UserPresence
		err = rows.Scan(&userPresence.Userathost, &userPresence.Presence, &userPresence.Lastpubtime)
		if err != nil {
			return userPresences, err
		}
		// and append it to the array
		userPresences = append(userPresences, userPresence)
	}
	return userPresences, err
}
