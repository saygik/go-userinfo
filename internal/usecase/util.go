package usecase

import (
	"encoding/json"
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
