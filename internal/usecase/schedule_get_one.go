package usecase

import (
	"fmt"
	"strconv"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSchedule(id int, currentUser string) (entity.Schedule, error) {

	schedule, err := u.repo.GetSchedule(id)
	if err != nil {
		return entity.Schedule{}, u.Error("невозможно получить календарь из базы данных")
	}
	schedule.Available = false
	scheduleAvailable := false
	scheduleUsersIdName, err := u.repo.GetScheduleUsers(id, 2)
	if err != nil {
		return entity.Schedule{}, u.Error("невозможно получить пользователей календаря из базы данных")
	}
	//GET GLPI Users of schedule
	glpiGroups, _ := u.repo.GetScheduleUserGroups(id, "glpi")

	for _, group := range glpiGroups {
		if group.Type != 2 {
			continue
		}
		gid, err := strconv.Atoi(group.Name)
		if err != nil {
			continue
		}
		groupUsers, _ := u.glpi.GetGroupUsers(gid)
		for _, groupUser := range groupUsers {
			if !IsStringInArrayIdName(groupUser.Name, scheduleUsersIdName) {
				scheduleUsersIdName = append(scheduleUsersIdName, groupUser)
			}
		}
	}

	scheduleUsers := []entity.ScheduleUser{}
	su := entity.ScheduleUser{}
	for _, user := range scheduleUsersIdName {
		if user.Name == currentUser {
			scheduleAvailable = true
		}
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

	//Get users with edit permission
	scheduleAdmins, err := u.repo.GetScheduleUsers(id, 1)
	if err != nil {
		return entity.Schedule{}, u.Error("ошибка получения списка доступа пользователей")
	}
	for _, group := range glpiGroups {
		if group.Type != 1 {
			continue
		}
		gid, err := strconv.Atoi(group.Name)
		if err != nil {
			continue
		}
		groupUsers, _ := u.glpi.GetGroupUsers(gid)
		for _, groupUser := range groupUsers {
			if !IsStringInArrayIdName(groupUser.Name, scheduleAdmins) {
				scheduleAdmins = append(scheduleAdmins, groupUser)
			}
		}
	}
	for _, user := range scheduleAdmins {
		if user.Name == currentUser {
			scheduleAvailable = true
			schedule.Edit = true
			break
		}
	}

	scheduleTasks, _ := u.GetScheduleTasks(id)
	schedule.ScheduleTasks = scheduleTasks
	schedule.ScheduleUsers = scheduleUsers
	schedule.ScheduleAdmins = scheduleAdmins
	if !schedule.Private {
		schedule.Available = true
	} else {
		schedule.Available = scheduleAvailable
	}

	//h.uc.GetADGroupUsers(domain, group)

	return schedule, nil
}
