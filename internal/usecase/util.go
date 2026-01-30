package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func CurrentTimeFormattedRFC3339() (string, error) {
	// Load the Europe/Minsk time zone
	loc, err := time.LoadLocation("Europe/Minsk")
	if err != nil {
		return "", errors.New("ошибка загрузки временной зоны Europe/Minsk")
	}
	tInMinsk := time.Now().In(loc)

	formatted := tInMinsk.Format(time.RFC3339)

	return formatted, nil
}

func IsStringInArray(str string, arr any) bool {
	if arr == nil {
		return false
	}
	return slices.Contains(arr.([]string), str)
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

func isStringObjsEqual(obj1 interface{}, obj2 interface{}) bool {
	str1, ok := obj1.(string)
	if !ok {
		return false
	}
	str2, ok := obj2.(string)
	if !ok {
		return false
	}
	return strings.ToUpper(str1) == strings.ToUpper(str2)
}

func ADFiletimeToGoTime(adFiletime string) (time.Time, error) {
	filetime := new(big.Int)
	if _, ok := filetime.SetString(adFiletime, 10); !ok {
		return time.Time{}, fmt.Errorf("invalid filetime string")
	}

	// Файлвремени - 100 наносекундные интервалы, переводим в миллисекунды
	msSince1601 := new(big.Int).Div(filetime, big.NewInt(10000))

	// Миллисекунды между 1601 и 1970 (эпоха unix)
	const msBetween1601And1970 int64 = 11644473600000

	// Вычисляем unix время в миллисекундах
	unixMs := new(big.Int).Sub(msSince1601, big.NewInt(msBetween1601And1970))

	unixMsInt64 := unixMs.Int64()

	// Конвертируем unixMs в time.Time
	return time.Unix(0, unixMsInt64*int64(time.Millisecond)), nil
}

// AnyOfFirstInSecond returns true if any element from 'first' exists in 'second' or if any element from 'second' starts with any element from 'first'.
func AnyOfFirstInSecond(first, second []string) bool {
	if len(first) == 0 || len(second) == 0 {
		return false
	}

	// Check exact matches
	lookup := make(map[string]struct{}, len(second))
	for _, v := range second {
		lookup[v] = struct{}{}
	}
	for _, v := range first {
		if _, ok := lookup[v]; ok {
			return true
		}
	}

	// Check if any string from 'second' starts with any string from 'first'
	for _, firstStr := range first {
		for _, secondStr := range second {
			if strings.HasPrefix(secondStr, firstStr) {
				return true
			}
		}
	}

	return false
}

func GetStringFromMap(m map[string]any, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
