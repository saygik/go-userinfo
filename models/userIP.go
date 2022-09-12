package models

import (
	"github.com/saygik/go-userinfo/db"
	"github.com/saygik/go-userinfo/forms"
)

// UserModel ...
type UserIPModel struct{}

type User struct {
	Login string `db:"login" json:"login"`
	Ip    string `db:"ip" json:"ip"`
}

type UserActivity struct {
	Ip       string `db:"ip" json:"ip"`
	Activity string `db:"activity" json:"activity"`
	Date     string `db:"date" json:"date"`
}

// GLPI User find by Mail ...
func (m UserIPModel) All(domain string) (users []User, err error) {
	_, err = db.GetDB().Select(&users, "GetAllUserIPByDomain $1", domain)
	return users, err
}
func (m UserIPModel) SetUserIp(form forms.UserActivityForm) (msgResponce string, err error) {
	//	_, err = db.GetDB().Exec("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy)
	err = db.GetDB().QueryRow("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy).Scan(&msgResponce)

	return msgResponce, err
}
func (m UserIPModel) GetUserByName(user string) (form forms.UserActivityForm, err error) {
	err = db.GetDB().QueryRow("GetUserByName $1", user).Scan(&form.User, &form.Ip, &form.Activiy, &form.ActiviyIp, &form.Date)
	return form, err
}
func (m UserIPModel) GetUserWeekActivity(login string) (activity []UserActivity, err error) {
	_, err = db.GetDB().Select(&activity, "GetUserLastWeekActivity $1", login)
	return activity, err
}

func (m UserIPModel) AllScheduleTasks(schedule string) (schedules []forms.ScheduleTask, err error) {
	_, err = db.GetDB().Select(&schedules, "GetScheduleTasks $1", schedule)
	return schedules, err
}
func (m UserIPModel) Schedule(id string) (schedule forms.Schedule, err error) {
	err = db.GetDB().SelectOne(&schedule, "GetSchedule $1", id)
	return schedule, err
}

func (m UserIPModel) AddScheduleTask(form forms.ScheduleTask) (msgResponce string, err error) {
	err = db.GetDB().QueryRow("AddScheduleTask $1,$2,$3,$4,$5", form.Idc, form.Title, form.Upn, form.Start, form.End).Scan(&msgResponce)

	return msgResponce, err
}

func (m UserIPModel) UpdateScheduleTask(form forms.ScheduleTask) (rows int64, err error) {
	res, err := db.GetDB().Exec("UpdateScheduleTask $1,$2,$3", form.Id, form.Start, form.End)
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return 0, err
		}
		return ra, err
	}
	return 0, err
}
func (m UserIPModel) DelScheduleTask(id int64) (rows int64, err error) {
	res, err := db.GetDB().Exec("DelScheduleTask $1", id)
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return 0, err
		}
		return ra, err
	}
	return 0, err
}
