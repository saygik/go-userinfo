package usecase

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSchedule(id int, currentUser string, isSysAdmin bool) (entity.Schedule, error) {
	schedule, err := u.repo.GetSchedule(id)
	if err != nil {
		return entity.Schedule{}, u.Error("невозможно получить календарь из базы данных")
	}

	userRows, groupRows, err := u.repo.GetScheduleAccessRows(id)
	if err != nil {
		return entity.Schedule{}, u.Error("невозможно получить права календаря из базы данных")
	}

	managers := []entity.IdName{}
	viewers := []entity.IdName{}

	for _, row := range userRows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			continue
		}
		user := entity.IdName{Name: name}
		switch row.Role {
		case 1:
			if !IsStringInArrayIdName(name, managers) {
				managers = append(managers, user)
			}
		default:
			if !IsStringInArrayIdName(name, viewers) {
				viewers = append(viewers, user)
			}
		}
	}

	for _, row := range groupRows {
		if row.GroupId <= 0 {
			continue
		}
		groupUsers, err := u.glpi.GetGroupUsers(row.GroupId)
		if err != nil {
			continue
		}
		for _, groupUser := range groupUsers {
			switch row.Role {
			case 1:
				if !IsStringInArrayIdName(groupUser.Name, managers) {
					managers = append(managers, groupUser)
				}
			default:
				if !IsStringInArrayIdName(groupUser.Name, viewers) {
					viewers = append(viewers, groupUser)
				}
			}
		}
	}

	isManager := IsStringInArrayIdName(currentUser, managers)
	isViewer := IsStringInArrayIdName(currentUser, viewers)

	schedule.Edit = isSysAdmin || isManager
	schedule.Available = isSysAdmin || !schedule.Private || isManager || isViewer

	// В список участников графика входят и полные администраторы, и
	// наблюдатели: полный доступ включает возможность быть назначенным.
	roster := make([]entity.IdName, 0, len(managers)+len(viewers))
	for _, user := range managers {
		if !IsStringInArrayIdName(user.Name, roster) {
			roster = append(roster, user)
		}
	}
	for _, user := range viewers {
		if !IsStringInArrayIdName(user.Name, roster) {
			roster = append(roster, user)
		}
	}

	scheduleUsers := make([]entity.ScheduleUser, 0, len(roster))
	for _, user := range roster {
		adUser := u.GetUserADPropertysShort(user.Name)
		su := entity.ScheduleUser{
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

	scheduleTasks, _ := u.GetScheduleTasks(id)

	schedule.ScheduleAdmins = managers
	schedule.ScheduleUsers = scheduleUsers
	schedule.ScheduleTasks = scheduleTasks

	return schedule, nil
}
