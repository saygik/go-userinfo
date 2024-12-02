package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSchedule(id string, currentUser string) (entity.Schedule, error) {

	schedule, err := u.repo.GetSchedule(id)
	if err != nil {
		return entity.Schedule{}, u.Error("имя пользователя в запросе отсутствует")
	}
	scheduleUsersIdName, err := u.repo.GetScheduleUsers(id, 2)
	if err != nil {
		return entity.Schedule{}, u.Error("имя пользователя в запросе отсутствует")
	}
	scheduleUsers := []entity.ScheduleUser{}
	su := entity.ScheduleUser{}
	for _, user := range scheduleUsersIdName {
		adUser := u.GetUserADPropertysShort(user.Name)
		su = entity.ScheduleUser{
			Name: fmt.Sprintf("%s", adUser["name"]),
			AD:   fmt.Sprintf("%t", adUser["findedInAD"]),
		}
		if _, ok := adUser["displayName"]; ok {
			su.DisplayName = fmt.Sprintf("%s", adUser["displayName"])
		}
		if _, ok := adUser["company"]; ok {
			su.Company = fmt.Sprintf("%s", adUser["company"])
		}
		if _, ok := adUser["title"]; ok {
			su.Title = fmt.Sprintf("%s", adUser["title"])
		}
		if _, ok := adUser["department"]; ok {
			su.Department = fmt.Sprintf("%s", adUser["department"])
		}
		if _, ok := adUser["mail"]; ok {
			su.Mail = fmt.Sprintf("%s", adUser["mail"])
		}
		if _, ok := adUser["telephoneNumber"]; ok {
			su.TelephoneNumber = fmt.Sprintf("%s", adUser["telephoneNumber"])
		}
		if _, ok := adUser["mobile"]; ok {
			su.Mobile = fmt.Sprintf("%s", adUser["mobile"])
		}
		scheduleUsers = append(scheduleUsers, su)
	}
	schedule.Edit = false
	scheduleAdmins, err := u.repo.GetScheduleUsers(id, 1)
	if err != nil {
		return entity.Schedule{}, u.Error("ошибка получения списка доступа пользователей")
	}
	for _, user := range scheduleAdmins {
		if user.Name == currentUser {
			schedule.Edit = true
			break
		}
	}
	scheduleTasks, _ := u.repo.GetScheduleTasks(id)
	schedule.ScheduleTasks = scheduleTasks
	schedule.ScheduleUsers = scheduleUsers
	schedule.ScheduleAdmins = scheduleAdmins

	//h.uc.GetADGroupUsers(domain, group)

	return schedule, nil
}
