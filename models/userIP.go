package models

import (
	"github.com/saygik/go-userinfo/db"
	"github.com/saygik/go-userinfo/forms"
)

// UserModel ...
type UserIPModel struct{}

type User struct {
	Login    string `db:"login" json:"login"`
	Ip       string `db:"ip" json:"ip"`
	Computer string `db:"computer" json:"computer"`
}
type SoftUser struct {
	Id        int64                  `db:"id" json:"id"`
	Name      string                 `db:"user_name" json:"name"`
	Login     string                 `db:"user_login" json:"login,omitempty"`
	Comment   string                 `db:"user_comment" json:"comment,omitempty"`
	Fio       string                 `db:"fio" json:"fio,omitempty"`
	External  bool                   `db:"external" json:"external,omitempty"`
	Propertys map[string]interface{} `json:"props"`
}

type UserActivity struct {
	Ip       string `db:"ip" json:"ip"`
	Activity string `db:"activity" json:"activity"`
	Date     string `db:"date" json:"date"`
}

type IdName struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type AppResources struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Edit string `db:"edit" json:"edit"`
}
type IdNameAvatar struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}

// GLPI User find by Mail ...
func (m UserIPModel) All(domain string) (users []User, err error) {
	_, err = db.GetDB().Select(&users, "GetAllUserIPByDomain $1", domain)
	return users, err
}
func (m UserIPModel) AllAvatars(domain string) (users []IdNameAvatar, err error) {
	_, err = db.GetDB().Select(&users, "GetAllUsersAvatars $1", domain)
	return users, err
}

/********** Avatars*/
func (m UserIPModel) GetUserAvatar(userID string) (avatar string, err error) {
	err = db.GetDB().QueryRow("GetUserAvatar $1", userID).Scan(&avatar)
	return avatar, err
}
func (m UserIPModel) SetUserAvatar(userID string, avatar string) (msgResponce string, err error) {

	//	err = db.GetDB().QueryRow("SetUserAvatar $1,$2", userID, avatar).Scan(&msgResponce)
	//	return msgResponce, err
	res, err := db.GetDB().Exec("SetUserAvatar $1,$2", userID, avatar)
	if res != nil {
		_, err1 := res.RowsAffected()
		if err1 != nil {
			return "No rows affected", err
		}
		return "Avatar updated or created", err
	}
	return "No rows affected", err

}

func (m UserIPModel) GetUserRoles(userID string) (roles []IdName, err error) {
	_, err = db.GetDB().Select(&roles, "GetUserRoles $1", userID)
	return roles, err
}
func (m UserIPModel) GetUserGroups(userID string) (groups []IdName, err error) {
	_, err = db.GetDB().Select(&groups, "GetUserGroups $1", userID)
	return groups, err
}
func (m UserIPModel) GetCurrentUserResources(userID string) (groups []AppResources, err error) {
	_, err = db.GetDB().Select(&groups, "GetUserResources $1", userID)
	return groups, err
}

func (m UserIPModel) GetUserResourceAccess(resouceID string, userID string) (access int, err error) {
	err = db.GetDB().QueryRow("GetUserResourceAccess $1,$2", resouceID, userID).Scan(&access)
	return access, err
}
func (m UserIPModel) GetUserSoftwares(userID string) (softwares []int64, err error) {
	_, err = db.GetDB().Select(&softwares, "GetUserSoftwares $1", userID)
	return softwares, err
}
func (m UserIPModel) GetSoftwareUsers(softID int64) (users []SoftUser, err error) {
	_, err = db.GetDB().Select(&users, "GetSoftwareUsers $1", softID)
	return users, err
}
func (m UserIPModel) SetUserIp(form forms.UserActivityForm) (msgResponce string, err error) {
	//	_, err = db.GetDB().Exec("SetUserIPActivity $1,$2,$3", form.User, form.Ip, form.Activiy)
	err = db.GetDB().QueryRow("SetUserIPActivityComputer $1,$2,$3,$4", form.User, form.Ip, form.Computer, form.Activiy).Scan(&msgResponce)

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

func (m UserIPModel) AddScheduleTask(form forms.ScheduleTask) (formRes forms.ScheduleTask, err error) {
	err = db.GetDB().QueryRow("AddScheduleTask $1,$2,$3,$4,$5,$6,$7,$8,$9,$10",
		form.Idc, form.Tip, form.Status, form.Title, form.Upn, form.Start, form.End, form.AllDay, form.SendMattermost, form.Comment).Scan(&formRes.Id,
		&formRes.Idc, &formRes.Tip, &formRes.Status, &formRes.Title, &formRes.Upn, &formRes.Start, &formRes.End, &formRes.AllDay, &formRes.SendMattermost, &formRes.Comment)

	return formRes, err
}

func (m UserIPModel) UpdateScheduleTask(form forms.ScheduleTask) (rows int64, err error) {
	res, err := db.GetDB().Exec("UpdateScheduleTask $1,$2,$3,$4,$5,$6,$7,$8,$9", form.Id, form.Tip, form.Status, form.Title, form.Start, form.End, form.AllDay, form.SendMattermost, form.Comment)
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

func (m UserIPModel) DelOneUserSoftware(user string, id int64) (rows int64, err error) {
	res, err := db.GetDB().Exec("DelOneUserSoftware $1, $2", user, id)
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return 0, err1
		}
		return ra, nil
	}
	return 0, err
}
func (m UserIPModel) AddOneUserSoftware(form forms.SoftwareForm) (rows int64, err error) {
	res, err := db.GetDB().Exec("AddOneUserSoftware $1,$2", form.User, form.Id)
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return 0, err1
		}
		return ra, nil
	}
	return 0, err
}
func (m UserIPModel) AddOneSoftwareUser(form forms.SoftwareUsersForm) (rows int64, err error) {
	res, err := db.GetDB().Exec("AddOneSoftwareUser $1,$2,$3,$4,$5,$6", form.User, form.Id, form.Login, form.Comment, form.Fio, form.External)
	if res != nil {
		ra, err1 := res.RowsAffected()
		if err1 != nil {
			return 0, err1
		}
		return ra, nil
	}
	return 0, err
}
