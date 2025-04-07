package usecase

import (
	"encoding/json"
)

func (u *UseCase) GetADUsers(user string) ([]map[string]interface{}, error) {

	domain := getDomainFromUserName(user)

	redisADUsers, err := u.redis.GetKeyFieldAll("ad")
	if err != nil {
		return nil, err
	}

	var res []map[string]interface{}
	for domainName, oneDomain := range redisADUsers {
		access := u.GetAccessToResource(domainName, user)
		if access == -1 && domain == domainName {
			access = 0
		}
		if access == -1 {
			continue
		}
		var r []map[string]interface{}
		json.Unmarshal([]byte(oneDomain), &r)
		accessToTechnicalInfo := access == 1
		accessToShortTechnicalInfo := access == 2
		for _, user := range r {
			delete(user, "employeeNumber")
			if !accessToTechnicalInfo && !accessToShortTechnicalInfo {
				delete(user, "ip")
				delete(user, "pwdLastSet")
				delete(user, "proxyAddresses")
				delete(user, "passwordDontExpire")
				delete(user, "passwordCantChange")
				delete(user, "distinguishedName")
				delete(user, "userAccountControl")
				delete(user, "memberOf")

				user["restricted"] = true
			}
			if !accessToTechnicalInfo && accessToShortTechnicalInfo {
				delete(user, "pwdLastSet")
				delete(user, "proxyAddresses")
				delete(user, "passwordDontExpire")
				delete(user, "passwordCantChange")
				delete(user, "distinguishedName")
				delete(user, "userAccountControl")
				delete(user, "memberOf")

				user["partial_restricted"] = true
			}
		}
		res = append(res, r...)
	}
	return res, nil
}
