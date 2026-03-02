package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetUserADPropertys(username string, perms entity.Permissions) (map[string]interface{}, error) {

	domain := getDomainFromUserName(username)

	accessLevel := u.GetAccessLevelForDomain(&perms, domain)
	if accessLevel == "none" {
		return nil, u.Error("нет прав на чтение свойств пользователя " + domain)
	}
	user, err := u.GetUser(username)
	if err != nil {
		return nil, err
	}

	if accessLevel == "user" {
		for _, field := range sensitiveFields {
			delete(user, field)
		}
	}

	return user, nil

}
