package usecase

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
