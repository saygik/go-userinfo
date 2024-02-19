package usecase

func (u *UseCase) GetUserADPropertys(username string, techUser string) (map[string]interface{}, error) {

	domain := getDomainFromUserName(username)
	domainName := getDomainFromUserName(techUser)

	access := u.GetAccessToResource(domain, techUser)
	if access == -1 && domain == domainName {
		access = 0
	}
	if access == -1 {
		return nil, u.Error("У вас недостаточно прав на просмотр данных пользователя")
	}

	accessToTechnicalInfo := access == 1

	user, err := u.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !accessToTechnicalInfo {
		delete(user, "ip")
		delete(user, "pwdLastSet")
		delete(user, "proxyAddresses")
		delete(user, "employeeNumber")
		delete(user, "passwordDontExpire")
		delete(user, "passwordCantChange")
		delete(user, "distinguishedName")
		delete(user, "userAccountControl")
		delete(user, "memberOf")

	}
	return user, nil

}
