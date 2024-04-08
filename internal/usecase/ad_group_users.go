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
