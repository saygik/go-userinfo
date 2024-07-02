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
