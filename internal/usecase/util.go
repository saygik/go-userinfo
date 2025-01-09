package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func IsStringInArray(str string, arr interface{}) bool {
	if arr == nil {
		return false
	}
	for _, b := range arr.([]string) {
		if b == str {
			return true
		}
	}
	return false
}

func IsStringInArrayIdName(str string, arr []entity.IdName) bool {
	if arr == nil {
		return false
	}
	for _, b := range arr {
		if b.Name == str {
			return true
		}
	}
	return false
}

func getDomainFromUserName(s string) string {
	domainArray := strings.Split(s, "@")
	if len(domainArray) < 2 {
		return ""
	}
	return domainArray[1]
}

func unmarshalString(str string, v any) error {
	return json.Unmarshal([]byte(str), v)
}

func isTechnicalAdminOfUser(user entity.GLPIUser, tech entity.GLPIUser) bool {

	if len(user.Profiles) < 1 {
		return false
	}
	if len(tech.Profiles) < 1 {
		return false
	}
	var userProfiles []int
	for _, tp := range tech.Profiles {
		for _, up := range user.Profiles {
			if tp.Id == 1 {
				continue
			}
			if tp.Recursive {
				if err := json.Unmarshal([]byte(up.Orgs), &userProfiles); err != nil {
					return false
				}
				if up.Eid == tp.Eid {
					return true
				}
				if containsInt(userProfiles, tp.Eid) {
					return true
				}
			} else {
				if up.Eid == tp.Eid {
					return true
				}
			}
		}
	}

	return false
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsIntInIdNameTypeArray(s []entity.IdNameType, itype int, e int) bool {
	for _, a := range s {
		if a.Id == e && (itype == 0 || a.Type == itype) {
			return true
		}
	}
	return false
}

func containsIDNameInIdNameTypeArray(s []entity.IdNameType, itype int, e []entity.IdName) bool {
	for _, a := range s {
		for _, b := range e {
			if a.Id == b.Id && (itype == 0 || a.Type == itype) {
				return true
			}
		}
	}
	return false
}

func filterStringArrayByWord(words []string, word string) []string {
	filtered := []string{}
	for i := range words {
		if strings.HasPrefix(words[i], word) {
			filtered = append(filtered, words[i])
		}
	}
	return filtered
}

func addStringToArrayIfNotExist(str string, arr []string) []string {
	if str == "" {
		return arr
	}
	for _, b := range arr {
		if b == str {
			return arr
		}
	}
	arr = append(arr, str)
	return arr
}

func (u *UseCase) scheduleTaskToScheduleTaskCalendar(task entity.ScheduleTask) entity.ScheduleTaskCalendar {

	oneTask := entity.ScheduleTaskCalendar{
		Id:     task.Id,
		Title:  task.Title,
		Start:  task.Start,
		End:    task.End,
		AllDay: task.AllDay}
	oneTask.ExtendedProps.Id = task.Upn
	oneTask.ExtendedProps.Comment = task.Comment
	oneTask.ExtendedProps.NotificationSended = task.NotificationSended
	oneTask.ExtendedProps.SendMattermost = task.SendMattermost
	oneTask.ExtendedProps.Status = task.Status
	oneTask.ExtendedProps.Tip = task.Tip
	oneTask.ExtendedProps.Title = task.Title
	if task.Tip == 1 {
		user, err := u.GetUserShort(task.Upn)
		if err != nil {
			oneTask.ExtendedProps.Notfound = true
		} else {
			oneTask.ExtendedProps.Notfound = false
			oneTask.ExtendedProps.Company = convertInterfaceToString(user["company"])
			oneTask.ExtendedProps.Department = convertInterfaceToString(user["department"])
			oneTask.ExtendedProps.Mail = convertInterfaceToString(user["mail"])
			oneTask.ExtendedProps.Mobile = convertInterfaceToString(user["mobile"])
			oneTask.ExtendedProps.TelephoneNumber = convertInterfaceToString(user["telephoneNumber"])
		}
	}

	return oneTask
}

func convertInterfaceToString(val interface{}) string {
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func parseHTML(s string) string {
	res := s
	res = strings.Replace(res, "&lt;span&gt;", "", -1)
	res = strings.Replace(res, "&lt;/span&gt;", "", -1)
	res = strings.Replace(res, "&lt;br&gt;", "\r\n", -1)
	res = strings.Replace(res, "&lt;br /&gt;", "\r\n", -1)
	res = strings.Replace(res, "&lt;p&gt;", "", -1)
	res = strings.Replace(res, "&lt;/p&gt;", "", -1)

	return res
}
