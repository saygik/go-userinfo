package usecase

import "fmt"

func (u *UseCase) GetADGroupUsers(domain string, group string) ([]map[string]interface{}, error) {

	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("такой домен отсутствует в системе")
	}

	users, err := u.ad.GetGroupUsers(domain, group)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка Active Directory: %s", err.Error()))
	}

	return users, err
}
func (u *UseCase) UserInDomainGroup(userID string, group string) error {
	domain := getDomainFromUserName(userID)

	if !u.ad.IsDomainExist(domain) {
		return u.Error("такой домен отсутствует в системе")
	}

	users, err := u.GetADGroupUsers(domain, group)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user["userPrincipalName"] == userID {
			return nil
		}
	}

	return u.Error("пользователь отсутствует в группе Active Directory")
}
func (u *UseCase) UserInDomainGroup2(userID string, group string, userDomain string) error {
	domain := getDomainFromUserName(userID)

	if !u.ad.IsDomainExist(domain) {
		return u.Error("такой домен отсутствует в системе")
	}
	if domain != userDomain {
		return u.Error("доступ в другие домены не предоставляется")
	}
	userinfo, err := u.ad.GetUser(domain, userID)
	if err != nil {
		return u.Error("Технический специалист не найден")
	}
	allGroups := []string{}
	if allGroupsTested, ok := userinfo["memberOf"].([]string); !ok {
		return u.Error("Отсутствуют группы у технического специалиста")
	} else {
		allGroups = allGroupsTested
	}

	for _, groupName := range allGroups {
		if groupName == group {
			return nil
		}
	}

	return u.Error("пользователь отсутствует в группе Active Directory")
}
