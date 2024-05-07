package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGlpiUserForTechnical(userName string, techName string) (*entity.GLPIUser, error) {
	emptyUser := entity.GLPIUser{}
	if userName == "" {
		return nil, u.Error("имя пользователя в запросе отсутствует")
	}
	if !isEmailValid(userName) {
		return &emptyUser, u.Error("неверное имя пользователя в запросе")
	}

	if techName == "" {
		return nil, u.Error("сначала войдите в систему")
	}

	userTech, err := u.GetGlpiUser(techName)
	if err != nil {
		return nil, u.Error("пользователь-техспециалист не найден в системе GLPI")
	}

	user, err := u.GetGlpiUser(userName)
	if err != nil {
		return nil, u.Error("пользователь не найден в системе GLPI")
	}

	GetTicketsNonClosedFromIniciator, err := u.glpi.GetTicketsNonClosedFromIniciator(user.Id)
	if err == nil {
		user.Tickets = GetTicketsNonClosedFromIniciator
	}
	if userName == techName {
		return &user, nil
	}
	if !isTechnicalAdminOfUser(user, userTech) {
		return nil, u.Error("у вас нет прав на этого пользователя в системе GLPI")
	}

	return &user, nil
}
