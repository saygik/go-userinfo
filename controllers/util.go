package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/saygik/go-userinfo/models"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	if strings.Index(IPAddress, "::1") > -1 {
		return "127.0.0.1"
	}
	i := strings.Index(IPAddress, ":")
	if i > -1 {
		return IPAddress[:i]
	} else {
		return IPAddress
	}

}

// func ReadUserName(ip string) (names []string, err error) {
// 	return net.LookupAddr(ip)
// }

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
func isTechnicalAdminOfUser(user models.GLPIUser, tech models.GLPIUser) bool {

	if len(user.Profiles) < 1 {
		return false
	}
	if len(tech.Profiles) < 1 {
		return false
	}
	var userProfiles []int64
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
				if containsInt64(userProfiles, tp.Eid) {
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

func containsInt64(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsInt64InIdNameTypeArray(s []models.IdNameType, e int64) bool {
	for _, a := range s {
		if a.Id == e {
			return true
		}
	}
	return false
}
func containsIDNameInIdNameTypeArray(s []models.IdNameType, e []models.IdName) bool {
	for _, a := range s {
		for _, b := range e {
			if a.Id == b.Id {
				return true
			}
		}
	}
	return false
}

func GetDomainFromUserName(s string) string {
	return strings.Split(fmt.Sprintf("%s", s), "@")[1]
}
