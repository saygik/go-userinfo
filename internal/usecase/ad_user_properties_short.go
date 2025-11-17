package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetUserADPropertysShort(username string) map[string]interface{} {

	user, err := u.GetUserShort(username)
	if err != nil {
		userProperties := map[string]interface{}{}
		userProperties["name"] = username
		userProperties["findedInAD"] = false
		return userProperties
	} else {
		delete(user, "ip")
		delete(user, "pwdLastSet")
		delete(user, "proxyAddresses")
		delete(user, "passwordDontExpire")
		delete(user, "passwordCantChange")
		delete(user, "distinguishedName")
		delete(user, "userAccountControl")
		delete(user, "memberOf")
		delete(user, "employeeNumber")
		delete(user, "presence")
		delete(user, "otherTelephone")
		user["name"] = username
		user["findedInAD"] = true
		return user
	}

}

func (u *UseCase) GetUserADPropertysSimple(username string) (*entity.SimpleUser, error) {
	userRes := entity.SimpleUser{Id: username}

	user, err := u.GetUserShort(username)
	if err != nil {
		return &userRes, err

	} else {
		if val, ok := user["displayName"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Name = deptStr
			}
		}
		if val, ok := user["company"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Company = deptStr
			}
		}
		if val, ok := user["department"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Department = deptStr
			}
		}
		if val, ok := user["title"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Title = deptStr
			}
		}
		if val, ok := user["mail"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Mail = deptStr
			}
		}

		if val, ok := user["telephoneNumber"]; ok {
			if deptStr, ok := val.(string); ok {
				userRes.Telephone = deptStr
			}
		}
		return &userRes, nil
	}
}
