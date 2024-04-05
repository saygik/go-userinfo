package usecase

import "encoding/json"

func (u *UseCase) GetADUsersPublicInfo(user string) ([]map[string]interface{}, error) {

	if !u.HasTechnicalRole(user) {
		return nil, u.Error("у вас нет прав на просмотр списка пользователей всех доменов")
	}
	redisADUsers, err := u.redis.GetKeyFieldAll("allusers")
	if err != nil {
		return nil, err
	}
	var users []map[string]interface{}
	for _, value := range redisADUsers {
		var user map[string]interface{}
		json.Unmarshal([]byte(value), &user)
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
		delete(user, "url")
		delete(user, "otherTelephone")
		user["findedInAD"] = true
		user["name"] = user["userPrincipalName"]
		users = append(users, user)
	}
	//	json.Unmarshal([]byte(redisADUsers), &users)
	return users, nil
}
