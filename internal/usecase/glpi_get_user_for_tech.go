package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGlpiUserForTechnical(userName string, techName string) (entity.GLPIUser, error) {
	if userName == "" {
		return entity.GLPIUser{}, u.Error("имя пользователя в запросе отсутствует")
	}
	if !isEmailValid(userName) {
		return entity.GLPIUser{}, u.Error("неверное имя пользователя в запросе")
	}

	if techName == "" {
		return entity.GLPIUser{}, u.Error("сначала войдите в систему")
	}

	userTech, err := u.GetGlpiUser(techName)
	if err != nil {
		return entity.GLPIUser{}, u.Error("пользователь-техспециалист не найден в системе GLPI")
	}

	user, err := u.GetGlpiUser(userName)
	if err != nil {
		return entity.GLPIUser{}, u.Error("пользователь не найден в системе GLPI")
	}
	if !isTechnicalAdminOfUser(user, userTech) {
		return entity.GLPIUser{}, u.Error("у вас нет прав на этого пользователя в системе GLPI")
	}

	return user, nil
}
