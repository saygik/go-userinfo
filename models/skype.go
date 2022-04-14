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

type ActiveConference struct {
	Id    string `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Users string `db:"users" json:"users"`
}

type ConferencePresence struct {
	Npp         string `db:"npp" json:"npp"`
	Id          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	ConfId      string `db:"confid" json:"confid"`
	UserName    string `db:"userName" json:"userName"`
	JoinTime    string `db:"joinTime" json:"joinTime"`
	Displayname string `db:"displayname" json:"displayname,omitempty"`
	Company     string `db:"company" json:"company,omitempty"`
	Department  string `db:"department" json:"department,omitempty"`
	Dolg        string `db:"dolg" json:"dolg,omitempty"`
}

//GLPI User find by Mail ...
func (m SkypeModel) AllPresences() (userPresences []UserPresence, err error) {
	_, err = db.GetDBSkype().Select(&userPresences, "Skype_GetAvailability")
	return userPresences, err
}
func (m SkypeModel) OnePresence(user string) (userPresence UserPresence, err error) {
	err = db.GetDBSkype().SelectOne(&userPresence, "Skype_GetOneUserAvailability $1", user)
	userPresence.Userathost = user
	return userPresence, err
}

//Skype Get Active Conferences ...
func (m SkypeModel) AllActiveConferences() (conferences []ActiveConference, err error) {
	_, err = db.GetDBSkype().Select(&conferences, "Skype_GetActiveConferences")
	return conferences, err
}

//Skype_GetConferencePresence 6452...
func (m SkypeModel) ConferencePresence(confID int64) (conferencePresence []ConferencePresence, err error) {
	_, err = db.GetDBSkype().Select(&conferencePresence, "Skype_GetConferencePresence $1", confID)
	return conferencePresence, err
}
